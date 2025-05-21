package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// measurementService implementa la lógica de negocio para mediciones
type measurementService struct {
	measurementRepo ports.IMeasurementRepository
	tagRepo         ports.ITagRepository
	recommendRepo   ports.IRecommendationRepository
}

// NewMeasurementService crea una nueva instancia de MeasurementService
func NewMeasurementService(
	measurementRepo ports.IMeasurementRepository,
	tagRepo ports.ITagRepository,
	recommendRepo ports.IRecommendationRepository,
) ports.IMeasurementService {
	return &measurementService{
		measurementRepo: measurementRepo,
		tagRepo:         tagRepo,
		recommendRepo:   recommendRepo,
	}
}

// Create crea una nueva medición
func (s *measurementService) Create(ctx context.Context, measurement *domain.Measurement) error {
	if err := measurement.Validate(); err != nil {
		return err
	}
	return s.measurementRepo.Create(ctx, measurement)
}

// GetByID obtiene una medición por su ID
func (s *measurementService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Measurement, error) {
	return s.measurementRepo.GetByID(ctx, id)
}

// GetByPatientID obtiene mediciones por ID de paciente
func (s *measurementService) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByPatientID(ctx, patientID)
}

// GetByUserID obtiene mediciones por ID de usuario
func (s *measurementService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByUserID(ctx, userID)
}

// GetByTagID obtiene mediciones por ID de etiqueta
func (s *measurementService) GetByTagID(ctx context.Context, tagID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByTagID(ctx, tagID)
}

// GetByRecommendationID obtiene mediciones por ID de recomendación
func (s *measurementService) GetByRecommendationID(ctx context.Context, recommendationID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByRecommendationID(ctx, recommendationID)
}

// GetByDateRange obtiene mediciones dentro de un rango de fechas
func (s *measurementService) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByDateRange(ctx, startDate, endDate)
}

// GetAll obtiene todas las mediciones
func (s *measurementService) GetAll(ctx context.Context) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetAll(ctx)
}

// Update actualiza una medición existente
func (s *measurementService) Update(ctx context.Context, measurement *domain.Measurement) error {
	if err := measurement.Validate(); err != nil {
		return err
	}
	return s.measurementRepo.Update(ctx, measurement)
}

// Delete elimina una medición por su ID
func (s *measurementService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.measurementRepo.Delete(ctx, id)
}

// AssignTag asigna una etiqueta a una medición
func (s *measurementService) AssignTag(ctx context.Context, measurementID, tagID uuid.UUID) error {
	// Verificar que la medición existe
	measurement, err := s.measurementRepo.GetByID(ctx, measurementID)
	if err != nil {
		return err
	}

	// Verificar que la etiqueta existe
	if tagID != uuid.Nil {
		_, err = s.tagRepo.GetByID(ctx, tagID)
		if err != nil {
			return err
		}
	}

	// Asignar la etiqueta
	measurement.SetTag(tagID)
	return s.measurementRepo.Update(ctx, measurement)
}

// AssignRecommendation asigna una recomendación a una medición
func (s *measurementService) AssignRecommendation(ctx context.Context, measurementID, recommendationID uuid.UUID) error {
	// Verificar que la medición existe
	measurement, err := s.measurementRepo.GetByID(ctx, measurementID)
	if err != nil {
		return err
	}

	// Verificar que la recomendación existe
	if recommendationID != uuid.Nil {
		_, err = s.recommendRepo.GetByID(ctx, recommendationID)
		if err != nil {
			return err
		}
	}

	// Asignar la recomendación
	measurement.SetRecommendation(recommendationID)
	return s.measurementRepo.Update(ctx, measurement)
}