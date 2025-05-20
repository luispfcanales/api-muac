package config

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
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
			Name:        "Administrador",
			Description: "Acceso completo al sistema",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Médico",
			Description: "Acceso a pacientes y mediciones",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Enfermero",
			Description: "Registro de mediciones",
			CreatedAt:   time.Now(),
		},
	}

	if err := tx.Create(&roles).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insertar localidades
	localities := []domain.Locality{
		{
			ID:          uuid.New(),
			Name:        "Hospital Central",
			Location:    "Ciudad Capital",
			Description: "Hospital principal",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Clínica Norte",
			Location:    "Zona Norte",
			Description: "Clínica comunitaria",
			CreatedAt:   time.Now(),
		},
	}

	if err := tx.Create(&localities).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insertar tags
	tags := []domain.Tag{
		{
			ID:          uuid.New(),
			Name:        "Urgente",
			Description: "Requiere atención inmediata",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Control",
			Description: "Medición de control rutinaria",
			CreatedAt:   time.Now(),
		},
	}

	if err := tx.Create(&tags).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insertar recomendaciones
	recommendations := []domain.Recommendation{
		{
			ID:                   uuid.New(),
			Name:                 "Desnutrición severa",
			Description:          "Requiere hospitalización",
			RecommendationUmbral: "<11.5cm",
			CreatedAt:            time.Now(),
		},
		{
			ID:                   uuid.New(),
			Name:                 "Desnutrición moderada",
			Description:          "Suplementación nutricional",
			RecommendationUmbral: "11.5-12.5cm",
			CreatedAt:            time.Now(),
		},
		{
			ID:                   uuid.New(),
			Name:                 "Normal",
			Description:          "Sin intervención requerida",
			RecommendationUmbral: ">12.5cm",
			CreatedAt:            time.Now(),
		},
	}

	if err := tx.Create(&recommendations).Error; err != nil {
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
