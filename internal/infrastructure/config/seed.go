package config

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase inserta datos iniciales en la base de datos
func SeedDatabase(db *gorm.DB) error {
	log.Println("Iniciando la siembra de datos iniciales...")

	// Verificar si ya existen datos en la tabla ROLE
	var roleCount int64
	if err := db.Model(&domain.Role{}).Count(&roleCount).Error; err != nil {
		return err
	}

	// Si ya hay roles, asumimos que los datos ya fueron sembrados
	if roleCount > 0 {
		log.Println("Los datos ya han sido sembrados anteriormente")
		return nil
	}

	// Iniciar una transacción para asegurar la integridad de los datos
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Función para hacer rollback en caso de error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insertar roles
	roles := []domain.Role{
		{
			ID:          uuid.New(),
			Name:        "ADMINISTRADOR",
			Description: "Acceso completo al sistema",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "SUPERVISOR",
			Description: "Acceso a pacientes y mediciones",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "APODERADO",
			Description: "Registro de mediciones",
			CreatedAt:   time.Now(),
		},
	}

	if err := tx.Create(&roles).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hashear la contraseña
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Crear usuarios con rol administrador
	users := []domain.User{
		{
			ID:           uuid.New(),
			Name:         "administrador",
			LastName:     "administrador",
			Username:     "administrador",
			Email:        "admin@example.com",
			DNI:          "12345678",
			PasswordHash: string(hashedPassword),
			RoleID:       roles[0].ID,
			CreatedAt:    time.Now(),
		},
	}

	if err := tx.Create(&users).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Println("Datos iniciales sembrados correctamente")
	return nil
}
