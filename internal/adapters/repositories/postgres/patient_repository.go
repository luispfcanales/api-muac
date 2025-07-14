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

// patientRepository implementa la interfaz IPatientRepository usando GORM
type patientRepository struct {
	db *gorm.DB
}

// NewPatientRepository crea una nueva instancia de PatientRepository
func NewPatientRepository(db *gorm.DB) ports.IPatientRepository {
	return &patientRepository{
		db: db,
	}
}

// Create inserta un nuevo paciente en la base de datos
func (r *patientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	result := r.db.WithContext(ctx).Create(patient)
	if result.Error != nil {
		return fmt.Errorf("error al crear paciente: %w", result.Error)
	}
	return nil
}

// GetByID obtiene un paciente por su ID
func (r *patientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	var patient domain.Patient
	result := r.db.WithContext(ctx).
		// Mediciones ordenadas por fecha (más recientes primero)
		Preload("Measurements", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Measurements.Tag").
		Preload("Measurements.Recommendation").
		Where("ID = ?", id).First(&patient)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPatientNotFound
		}
		return nil, fmt.Errorf("error al obtener paciente: %w", result.Error)
	}
	return &patient, nil
}

// GetByDNI obtiene un paciente por su DNI
func (r *patientRepository) GetByDNI(ctx context.Context, dni string) (*domain.Patient, error) {
	var patient domain.Patient
	result := r.db.WithContext(ctx).Where("DNI = ?", dni).First(&patient)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPatientNotFound
		}
		return nil, fmt.Errorf("error al obtener paciente por DNI: %w", result.Error)
	}
	return &patient, nil
}

// GetAll obtiene todos los pacientes
func (r *patientRepository) GetAll(ctx context.Context) ([]*domain.Patient, error) {
	var patients []*domain.Patient
	result := r.db.WithContext(ctx).Find(&patients)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener pacientes: %w", result.Error)
	}
	return patients, nil
}

// Update actualiza un paciente existente
func (r *patientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	result := r.db.WithContext(ctx).Save(patient)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar paciente: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrPatientNotFound
	}
	return nil
}

// Delete elimina un paciente por su ID
func (r *patientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Patient{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar paciente: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrPatientNotFound
	}
	return nil
}

// GetByFatherID obtiene los pacientes asociados a un padre específico
func (r *patientRepository) GetByFatherID(ctx context.Context, fatherID uuid.UUID) ([]*domain.Patient, error) {
	var patients []*domain.Patient
	// Asumiendo que hay una tabla de relación entre Father y Patient
	result := r.db.WithContext(ctx).
		Joins("JOIN FATHER ON FATHER.PATIENT_ID = PATIENT.ID").
		Where("FATHER.ID = ?", fatherID).
		Find(&patients)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener pacientes por ID de padre: %w", result.Error)
	}
	return patients, nil
}

// GetMeasurements obtiene todas las mediciones de un paciente específico
func (r *patientRepository) GetMeasurements(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error) {
	var measurements []*domain.Measurement
	result := r.db.WithContext(ctx).
		Where("PATIENT_ID = ?", patientID).
		Find(&measurements)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener mediciones del paciente: %w", result.Error)
	}
	return measurements, nil
}

// GetPatientsInRisk obtiene todos los pacientes en riesgo con todos sus datos
// Pacientes en riesgo = MUAC < 12.5 cm (moderado y severo)
func (r *patientRepository) GetPatientsInRisk(ctx context.Context, filters *domain.ReportFilters) ([]*domain.Patient, error) {
	var patients []*domain.Patient

	// Subconsulta para obtener la medición más reciente de cada paciente
	subQuery := r.db.Table("measurements m1").
		Select("m1.patient_id, MAX(m1.created_at) as max_created_at").
		Group("m1.patient_id")

	// Consulta principal
	query := r.db.WithContext(ctx).
		Table("patients p").
		Joins("JOIN measurements m ON p.id = m.patient_id").
		Joins("JOIN (?) latest ON m.patient_id = latest.patient_id AND m.created_at = latest.max_created_at", subQuery).
		Where("m.muac_value < ?", 12.5). // Solo pacientes en riesgo (< 12.5 cm)
		Preload("Measurements", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1) // Solo la última medición
		}).
		Preload("Measurements.Tag").
		Preload("Measurements.Recommendation").
		Order("m.muac_value ASC") // Los más críticos primero

	// Aplicar filtros si existen
	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Joins("JOIN users u ON p.user_id = u.id").
				Where("u.locality_id = ?", *filters.LocalityID)
		}
		if filters.Days > 0 {
			since := time.Now().AddDate(0, 0, -filters.Days)
			query = query.Where("m.created_at >= ?", since)
		}
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit)
		}
	}

	result := query.Find(&patients)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener pacientes en riesgo: %w", result.Error)
	}

	return patients, nil
}
