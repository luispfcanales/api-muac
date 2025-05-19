package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
	_ "github.com/lib/pq"              // Driver para PostgreSQL
	_ "github.com/mattn/go-sqlite3"    // Driver para SQLite
)

// DBType representa el tipo de base de datos
type DBType string

const (
	// PostgreSQL representa una base de datos PostgreSQL
	PostgreSQL DBType = "postgres"
	// MySQL representa una base de datos MySQL
	MySQL DBType = "mysql"
	// SQLite representa una base de datos SQLite
	SQLite DBType = "sqlite"
)

// Config contiene la configuración de la aplicación
type Config struct {
	// Tipo de base de datos (postgres, mysql, sqlite)
	DBType     DBType
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	// Ruta del archivo para SQLite
	DBPath string
	// Ruta del archivo SQL para inicializar la base de datos
	SQLFilePath string
	ServerPort  int
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	dbType := DBType(getEnv("DB_TYPE", string(SQLite))) // Por defecto SQLite

	return &Config{
		DBType:      dbType,
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      dbPort,
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "MUAC"),
		DBPath:      getEnv("DB_PATH", "./muac.db"),
		SQLFilePath: getEnv("SQL_FILE_PATH", "./ddbb.sql"),
		ServerPort:  serverPort,
	}
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// NewDBConnection crea una nueva conexión a la base de datos
func NewDBConnection(config *Config) (*sql.DB, error) {
	var dsn string
	var driverName string

	switch config.DBType {
	case PostgreSQL:
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	case MySQL:
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	case SQLite:
		driverName = "sqlite3"
		dsn = config.DBPath
	default:
		return nil, fmt.Errorf("tipo de base de datos no soportado: %s", config.DBType)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Si es SQLite, inicializar la base de datos si es necesario
	if config.DBType == SQLite {
		if err := InitializeSQLiteDBFromFile(db, config.SQLFilePath); err != nil {
			return nil, fmt.Errorf("error al inicializar la base de datos SQLite: %w", err)
		}
	}

	return db, nil
}

// InitializeSQLiteDBFromFile inicializa la base de datos SQLite con las tablas y datos del archivo SQL
func InitializeSQLiteDBFromFile(db *sql.DB, sqlFilePath string) error {
	// Verificar si las tablas ya existen
	var tableCount int
	err := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='ROLE'").Scan(&tableCount)
	if err != nil {
		return err
	}

	// Si la tabla ROLE ya existe, asumimos que la base de datos ya está inicializada
	if tableCount > 0 {
		return nil
	}

	// Leer el archivo SQL
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo SQL: %w", err)
	}

	sqlContent := string(sqlBytes)

	// Adaptar el SQL para SQLite
	sqlContent = adaptSQLForSQLite(sqlContent)

	// Dividir el contenido en sentencias SQL individuales
	statements := splitSQLStatements(sqlContent)

	// Iniciar una transacción para ejecutar todas las sentencias
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Ejecutar cada sentencia SQL
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error al ejecutar sentencia SQL: %s, error: %w", stmt, err)
		}
	}

	// Confirmar la transacción
	return tx.Commit()
}

// adaptSQLForSQLite adapta el SQL para que sea compatible con SQLite
func adaptSQLForSQLite(sql string) string {
	// Reemplazar UUID por TEXT
	sql = strings.ReplaceAll(sql, "UUID PRIMARY KEY", "TEXT PRIMARY KEY")

	// Reemplazar funciones específicas de PostgreSQL
	sql = strings.ReplaceAll(sql, "uuid_generate_v4()", "lower(hex(randomblob(16)))")
	sql = strings.ReplaceAll(sql, "CURRENT_DATE", "date('now')")

	// Eliminar comandos específicos de PostgreSQL/MySQL
	sql = removeLines(sql, "CREATE DATABASE")
	sql = removeLines(sql, "USE ")

	return sql
}

// removeLines elimina líneas que contienen cierto texto
func removeLines(input, contains string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if !strings.Contains(line, contains) {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// splitSQLStatements divide el contenido SQL en sentencias individuales
func splitSQLStatements(sql string) []string {
	// Dividir por punto y coma, pero ignorar los que están dentro de comillas
	var statements []string
	var currentStmt strings.Builder
	inQuote := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(sql); i++ {
		char := sql[i]

		// Manejar comentarios de línea
		if i < len(sql)-1 && char == '-' && sql[i+1] == '-' && !inQuote && !inBlockComment {
			inLineComment = true
			currentStmt.WriteByte(char)
			continue
		}

		// Manejar fin de comentarios de línea
		if (char == '\n' || char == '\r') && inLineComment {
			inLineComment = false
			currentStmt.WriteByte(char)
			continue
		}

		// Manejar comentarios de bloque
		if i < len(sql)-1 && char == '/' && sql[i+1] == '*' && !inQuote && !inLineComment {
			inBlockComment = true
			currentStmt.WriteByte(char)
			continue
		}

		// Manejar fin de comentarios de bloque
		if i < len(sql)-1 && char == '*' && sql[i+1] == '/' && inBlockComment {
			inBlockComment = false
			currentStmt.WriteByte(char)
			currentStmt.WriteByte(sql[i+1])
			i++
			continue
		}

		// Manejar comillas
		if char == '\'' && !inLineComment && !inBlockComment {
			inQuote = !inQuote
		}

		// Manejar punto y coma fuera de comillas y comentarios
		if char == ';' && !inQuote && !inLineComment && !inBlockComment {
			currentStmt.WriteByte(char)
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
			continue
		}

		// Agregar el carácter actual a la sentencia
		currentStmt.WriteByte(char)
	}

	// Agregar la última sentencia si no está vacía
	lastStmt := strings.TrimSpace(currentStmt.String())
	if lastStmt != "" {
		statements = append(statements, lastStmt)
	}

	return statements
}
