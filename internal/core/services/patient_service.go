package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// patientService implementa la lógica de negocio para pacientes
type patientService struct {
	patientRepo     ports.IPatientRepository
	measurementRepo ports.IMeasurementRepository
}

// NewPatientService crea una nueva instancia de PatientService
func NewPatientService(patientRepo ports.IPatientRepository, measurementRepo ports.IMeasurementRepository) ports.IPatientService {
	return &patientService{
		patientRepo:     patientRepo,
		measurementRepo: measurementRepo,
	}
}

// Create crea un nuevo paciente
func (s *patientService) Create(ctx context.Context, patient *domain.Patient) error {
	if err := patient.Validate(); err != nil {
		return err
	}
	//validar que no se repita el dni con otro registro
	p, err := s.patientRepo.GetByDNI(ctx, patient.DNI)
	if err != nil {
		return err
	}
	if p.ID != uuid.Nil {
		return errors.New("el DNI ya está registrado")
	}
	return s.patientRepo.Create(ctx, patient)
}

// GetByID obtiene un paciente por su ID
func (s *patientService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	return s.patientRepo.GetByID(ctx, id)
}

// GetByDNI obtiene un paciente por su DNI
func (s *patientService) GetByDNI(ctx context.Context, dni string) (*domain.Patient, error) {
	return s.patientRepo.GetByDNI(ctx, dni)
}

// GetAll obtiene todos los pacientes
func (s *patientService) GetAll(ctx context.Context) ([]*domain.Patient, error) {
	return s.patientRepo.GetAll(ctx)
}

// Update actualiza un paciente existente
func (s *patientService) Update(ctx context.Context, patient *domain.Patient) error {
	if err := patient.Validate(); err != nil {
		return err
	}
	return s.patientRepo.Update(ctx, patient)
}

// Delete elimina un paciente por su ID
func (s *patientService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.patientRepo.Delete(ctx, id)
}

// GetByFatherID obtiene los pacientes asociados a un padre específico
func (s *patientService) GetByFatherID(ctx context.Context, fatherID uuid.UUID) ([]*domain.Patient, error) {
	return s.patientRepo.GetByFatherID(ctx, fatherID)
}

// GetMeasurements obtiene todas las mediciones de un paciente específico
func (s *patientService) GetMeasurements(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error) {
	return s.patientRepo.GetMeasurements(ctx, patientID)
}

// AddMeasurement añade una nueva medición a un paciente
func (s *patientService) AddMeasurement(ctx context.Context, patientID uuid.UUID, measurement *domain.Measurement) error {
	// Verificar que el paciente existe
	_, err := s.patientRepo.GetByID(ctx, patientID)
	if err != nil {
		return err
	}

	// Asignar el ID del paciente a la medición
	measurement.PatientID = patientID

	// Validar la medición
	if err := measurement.Validate(); err != nil {
		return err
	}

	// Guardar la medición en la base de datos usando el repositorio de mediciones
	return s.measurementRepo.Create(ctx, measurement)
}
