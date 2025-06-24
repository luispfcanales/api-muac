package config

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
	_ "github.com/lib/pq"              // Driver para PostgreSQL
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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
	ServerPort int
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8003"))
	dbType := DBType(getEnv("DB_TYPE", string(PostgreSQL)))

	return &Config{
		DBType:     dbType,
		DBHost:     getEnv("DB_HOST", "35.173.114.173"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "unamadconfericis"),
		DBPassword: getEnv("DB_PASSWORD", "unamad2024."),
		DBName:     getEnv("DB_NAME", "muac"),
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
	default:
		return nil, fmt.Errorf("tipo de base de datos no soportado: %s", config.DBType)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
