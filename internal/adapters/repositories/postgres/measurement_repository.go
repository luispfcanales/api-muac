package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// measurementRepository implementa la interfaz IMeasurementRepository usando GORM
type measurementRepository struct {
	db *gorm.DB
}

// NewMeasurementRepository crea una nueva instancia de MeasurementRepository
func NewMeasurementRepository(db *gorm.DB) ports.IMeasurementRepository {
	return &measurementRepository{
		db: db,
	}
}

// Create inserta una nueva medición en la base de datos
func (r *measurementRepository) Create(ctx context.Context, measurement *domain.Measurement) error {
	result := r.db.WithContext(ctx).Create(measurement)
	if result.Error != nil {
		return fmt.Errorf("error al crear medición: %w", result.Error)
	}
	return nil
}

// GetByID obtiene una medición por su ID
func (r *measurementRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Measurement, error) {
	var measurement domain.Measurement
	result := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").
		Where("ID = ?", id).
		First(&measurement)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrMeasurementNotFound
		}
		return nil, fmt.Errorf("error al obtener medición: %w", result.Error)
	}
	return &measurement, nil
}

// GetByPatientID obtiene mediciones por ID de paciente
func (r *measurementRepository) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement
	result := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").
		Where("PATIENT_ID = ?", patientID).
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones por ID de paciente: %w", result.Error)
	}
	return measurements, nil
}

// GetByUserID obtiene mediciones por ID de usuario
func (r *measurementRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement
	result := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").
		Where("USER_ID = ?", userID).
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones por ID de usuario: %w", result.Error)
	}
	return measurements, nil
}

// GetByTagID obtiene mediciones por ID de etiqueta
func (r *measurementRepository) GetByTagID(ctx context.Context, tagID uuid.UUID) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement
	result := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").
		Where("TAG_ID = ?", tagID).
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones por ID de etiqueta: %w", result.Error)
	}
	return measurements, nil
}

// GetByRecommendationID obtiene mediciones por ID de recomendación
func (r *measurementRepository) GetByRecommendationID(ctx context.Context, recommendationID uuid.UUID) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement
	result := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").
		Where("RECOMMENDATION_ID = ?", recommendationID).
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones por ID de recomendación: %w", result.Error)
	}
	return measurements, nil
}

// GetByDateRange obtiene mediciones dentro de un rango de fechas
func (r *measurementRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement
	result := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").
		Where("TIMESTAMP BETWEEN ? AND ?", startDate, endDate).
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones por rango de fechas: %w", result.Error)
	}
	return measurements, nil
}

// GetAll obtiene todas las mediciones con todas sus relaciones ordenadas
func (r *measurementRepository) GetAll(ctx context.Context) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement

	result := r.db.WithContext(ctx).
		// Relaciones principales de Measurement
		Preload("Patient").
		Preload("User").
		Preload("Tag").
		Preload("Recommendation").

		// Relaciones anidadas del Patient
		Preload("Patient.User").                        // Usuario que creó el paciente
		Preload("Patient.User.Role").                   // Rol del usuario
		Preload("Patient.User.Locality").               // Localidad del usuario
		Preload("Patient.Measurements").                // Otras mediciones del paciente
		Preload("Patient.Measurements.Tag").            // Tags de otras mediciones
		Preload("Patient.Measurements.Recommendation"). // Recomendaciones de otras mediciones

		// Relaciones del User (quien tomó la medición)
		Preload("User.Role").     // Rol del usuario que midió
		Preload("User.Locality"). // Localidad del usuario que midió
		Preload("User.Patients"). // Pacientes asignados al usuario

		// Ordenamiento: más recientes primero
		Order("created_at DESC").
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones: %w", result.Error)
	}

	return measurements, nil
}

// Update actualiza una medición existente
func (r *measurementRepository) Update(ctx context.Context, measurement *domain.Measurement) error {
	result := r.db.WithContext(ctx).Save(measurement)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar medición: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrMeasurementNotFound
	}
	return nil
}

// Delete elimina una medición por su ID
func (r *measurementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Measurement{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar medición: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrMeasurementNotFound
	}
	return nil
}
