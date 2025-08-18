package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// patientService implementa la lógica de negocio para pacientes
type patientService struct {
	patientRepo     ports.IPatientRepository
	measurementRepo ports.IMeasurementRepository
	tipService      ports.ITipService
	recipeService   ports.IRecipeService
}

// NewPatientService crea una nueva instancia de PatientService
func NewPatientService(
	patientRepo ports.IPatientRepository,
	measurementRepo ports.IMeasurementRepository,
	tipService ports.ITipService,
	recipeService ports.IRecipeService,
) ports.IPatientService {
	return &patientService{
		patientRepo:     patientRepo,
		measurementRepo: measurementRepo,
		tipService:      tipService,
		recipeService:   recipeService,
	}
}

// Create crea un nuevo paciente
func (s *patientService) Create(ctx context.Context, patient *domain.Patient) error {
	if err := patient.Validate(); err != nil {
		return err
	}
	//validar que no se repita el dni con otro registro
	_, err := s.patientRepo.GetByDNI(ctx, patient.DNI)
	if err != nil {
		return s.patientRepo.Create(ctx, patient)
	}
	return domain.ErrPatientDNIAlreadyExists
}

// GetByID obtiene un paciente por su ID
func (s *patientService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	return s.patientRepo.GetByID(ctx, id)
}

// GetByDNI obtiene un paciente por su DNI
func (s *patientService) GetByDNI(ctx context.Context, dni string) (*domain.Patient, error) {
	patient, err := s.patientRepo.GetByDNI(ctx, dni)
	if err != nil {
		return nil, err
	}

	for i := range patient.Measurements {
		// Obtener tips y recetas para esta medición

		tips, _ := s.tipService.List(ctx, patient.Measurements[i].Tag.MuacCode)
		recipes, _ := s.recipeService.ListRecipesByAge(ctx, patient.Age)

		patient.Measurements[i].MeasurementAdvice = domain.MeasurementAdvice{
			Tips:    tips,
			Recipes: recipes,
		}
	}

	return patient, nil
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

// GetUsersWithRiskPatients obtiene usuarios con pacientes en riesgo
func (s *patientService) GetUsersWithRiskPatients(ctx context.Context, filters *domain.ReportFilters) ([]*domain.User, error) {
	// if err := s.ValidateFilters(filters); err != nil {
	// 	return nil, err
	// }

	users, err := s.patientRepo.GetUsersWithRiskPatients(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios con pacientes en riesgo: %w", err)
	}

	return users, nil
}
