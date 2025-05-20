package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
	_ "github.com/lib/pq"              // Driver para PostgreSQL
	_ "github.com/mattn/go-sqlite3"    // Driver para SQLite
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

// adaptSQLForSQLite adapta el SQL para que sea compatible con SQLite
func adaptSQLForSQLite(sql string) string {
	// Reemplazar UUID por TEXT
	sql = strings.ReplaceAll(sql, "UUID PRIMARY KEY", "TEXT PRIMARY KEY")

	// Reemplazar funciones específicas de PostgreSQL
	sql = strings.ReplaceAll(sql, "uuid_generate_v4()",
		"substr(lower(hex(randomblob(4))),1,8) || '-' || "+
			"substr(lower(hex(randomblob(2))),1,4) || '-' || "+
			"substr(lower(hex(randomblob(2))),1,4) || '-' || "+
			"substr(lower(hex(randomblob(2))),1,4) || '-' || "+
			"substr(lower(hex(randomblob(6))),1,12)")
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

	// Primero, eliminar todos los comentarios de línea
	lines := strings.Split(sql, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Ignorar líneas de comentario completas
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}
		// Para líneas con comentarios parciales, mantener solo la parte antes del comentario
		commentIndex := strings.Index(line, "--")
		if commentIndex >= 0 {
			line = line[:commentIndex]
		}
		if len(strings.TrimSpace(line)) > 0 {
			cleanedLines = append(cleanedLines, line)
		}
	}

	// Reconstruir el SQL sin comentarios de línea
	cleanedSQL := strings.Join(cleanedLines, " ")

	// Eliminar comentarios de bloque
	for {
		startIdx := strings.Index(cleanedSQL, "/*")
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(cleanedSQL[startIdx:], "*/")
		if endIdx == -1 {
			break
		}

		cleanedSQL = cleanedSQL[:startIdx] + " " + cleanedSQL[startIdx+endIdx+2:]
	}

	// Normalizar espacios en blanco
	cleanedSQL = strings.Join(strings.Fields(cleanedSQL), " ")

	// Dividir por punto y coma
	rawStatements := strings.Split(cleanedSQL, ";")

	// Filtrar sentencias vacías
	for _, stmt := range rawStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

// NewGormDBConnection crea una nueva conexión a la base de datos usando GORM
func NewGormDBConnection(config *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch config.DBType {
	case PostgreSQL:
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case MySQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case SQLite:
		db, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	default:
		return nil, fmt.Errorf("tipo de base de datos no soportado: %s", config.DBType)
	}

	if err != nil {
		return nil, err
	}

	// Si es SQLite y el archivo no existe, inicializar la base de datos
	if config.DBType == SQLite {
		if _, err := os.Stat(config.DBPath); os.IsNotExist(err) {
			if err := InitializeGormDBFromFile(db, config.SQLFilePath); err != nil {
				return nil, fmt.Errorf("error al inicializar la base de datos SQLite: %w", err)
			}
		}
	}

	return db, nil
}

// InitializeGormDBFromFile inicializa la base de datos con GORM
func InitializeGormDBFromFile(db *gorm.DB, sqlFilePath string) error {
	// Leer el archivo SQL
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo SQL: %w", err)
	}

	sqlContent := string(sqlBytes)

	// Adaptar el SQL para SQLite si es necesario
	sqlContent = adaptSQLForSQLite(sqlContent)

	// Dividir el contenido en sentencias SQL individuales
	statements := splitSQLStatements(sqlContent)

	// Ejecutar cada sentencia SQL
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("error al ejecutar sentencia SQL: %s, error: %w", stmt, err)
		}
	}

	return nil
}
