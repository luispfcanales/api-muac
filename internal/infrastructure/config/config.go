package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

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
	DBPath     string
	ServerPort int
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	dbType := DBType(getEnv("DB_TYPE", string(SQLite))) // Por defecto SQLite

	return &Config{
		DBType:     dbType,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "MUAC"),
		DBPath:     getEnv("DB_PATH", "./muac.db"),
		ServerPort: serverPort,
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

	return db, nil
}