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

// GetUsersWithRiskPatientsSimple - Versión corregida para UserID como puntero
func (r *patientRepository) GetUsersWithRiskPatients(ctx context.Context, filters *domain.ReportFilters) ([]*domain.User, error) {
	var users []*domain.User

	// PASO 1: Obtener todos los usuarios con sus pacientes y mediciones
	query := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Locality").
		Preload("Patients").
		Preload("Patients.Measurements", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC") // TODAS las mediciones, luego filtraremos en memoria
		}).
		Preload("Patients.Measurements.Tag").
		Preload("Patients.Measurements.Recommendation")

	// Aplicar filtros de usuario
	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Where("locality_id = ?", *filters.LocalityID)
		}
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit * 2) // Multiplicar por 2 para asegurar suficientes resultados
		}
	}

	var allUsers []*domain.User
	if err := query.Find(&allUsers).Error; err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}

	// PASO 2: Filtrar en memoria
	for _, user := range allUsers {
		var riskPatients []domain.Patient

		for _, patient := range user.Patients {
			if len(patient.Measurements) > 0 {
				// Tomar solo la última medición
				lastMeasurement := patient.Measurements[0]

				// Aplicar filtro de días si existe
				if filters != nil && filters.Days > 0 {
					since := time.Now().AddDate(0, 0, -filters.Days)
					if lastMeasurement.CreatedAt.Before(since) {
						continue
					}
				}

				// Verificar si está en riesgo
				if lastMeasurement.MuacValue < 12.5 {
					// Crear copia del paciente con solo la última medición
					riskPatient := patient
					riskPatient.Measurements = []domain.Measurement{lastMeasurement} // Solo la última medición
					riskPatients = append(riskPatients, riskPatient)
				}
			}
		}

		// Solo agregar usuario si tiene pacientes en riesgo
		if len(riskPatients) > 0 {
			user.Patients = riskPatients
			users = append(users, user)
		}
	}

	return users, nil
}

// GetPatientsInRisk obtiene todos los pacientes en riesgo con todos sus datos - CORREGIDO
// func (r *patientRepository) GetPatientsInRisk(ctx context.Context, filters *domain.ReportFilters) ([]*domain.Patient, error) {
// 	var patients []*domain.Patient

// 	// PASO 1: Obtener IDs de pacientes que están en riesgo según su última medición
// 	var patientIDs []uuid.UUID

// 	subQuery := r.db.Table("measurements m1").
// 		Select("m1.patient_id, MAX(m1.created_at) as max_created_at").
// 		Group("m1.patient_id")

// 	// Obtener solo los IDs de pacientes en riesgo
// 	result := r.db.WithContext(ctx).
// 		Table("patients p").
// 		Select("DISTINCT p.id").
// 		Joins("JOIN measurements m ON p.id = m.patient_id").
// 		Joins("JOIN (?) latest ON m.patient_id = latest.patient_id AND m.created_at = latest.max_created_at", subQuery).
// 		Where("m.muac_value < ?", 12.5)

// 	// Aplicar filtros si existen
// 	if filters != nil {
// 		if filters.LocalityID != nil {
// 			result = result.Joins("JOIN users u ON p.user_id = u.id").
// 				Where("u.locality_id = ?", *filters.LocalityID)
// 		}
// 		if filters.Days > 0 {
// 			since := time.Now().AddDate(0, 0, -filters.Days)
// 			result = result.Where("m.created_at >= ?", since)
// 		}
// 	}

// 	if err := result.Pluck("id", &patientIDs).Error; err != nil {
// 		return nil, fmt.Errorf("error al obtener IDs de pacientes en riesgo: %w", err)
// 	}

// 	// Si no hay pacientes en riesgo, retornar vacío
// 	if len(patientIDs) == 0 {
// 		return patients, nil
// 	}

// 	// PASO 2: Obtener pacientes completos con sus mediciones
// 	query := r.db.WithContext(ctx).
// 		Where("id IN ?", patientIDs).
// 		Preload("Measurements", func(db *gorm.DB) *gorm.DB {
// 			return db.Order("created_at DESC") // Todas las mediciones ordenadas por fecha
// 		}).
// 		Preload("Measurements.Tag").
// 		Preload("Measurements.Recommendation")

// 	// Aplicar límite si existe
// 	if filters != nil && filters.Limit > 0 {
// 		query = query.Limit(filters.Limit)
// 	}

// 	// Ordenar por el valor MUAC de la última medición (los más críticos primero)
// 	// Para esto necesitamos un ORDER BY especial
// 	query = query.Joins(`
// 		JOIN (
// 			SELECT m.patient_id, m.muac_value
// 			FROM measurements m
// 			JOIN (
// 				SELECT patient_id, MAX(created_at) as max_created_at
// 				FROM measurements
// 				GROUP BY patient_id
// 			) latest ON m.patient_id = latest.patient_id AND m.created_at = latest.max_created_at
// 		) last_muac ON patients.id = last_muac.patient_id
// 	`).Order("last_muac.muac_value ASC")

// 	if err := query.Find(&patients).Error; err != nil {
// 		return nil, fmt.Errorf("error al obtener pacientes completos: %w", err)
// 	}

// 	return patients, nil
// }
