package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// recommendationService implementa la lógica de negocio para recomendaciones
type recommendationService struct {
	recommendationRepo ports.IRecommendationRepository
}

// NewRecommendationService crea una nueva instancia de RecommendationService
func NewRecommendationService(recommendationRepo ports.IRecommendationRepository) ports.IRecommendationService {
	return &recommendationService{
		recommendationRepo: recommendationRepo,
	}
}

// Create crea una nueva recomendación
func (s *recommendationService) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	if err := recommendation.Validate(); err != nil {
		return err
	}
	return s.recommendationRepo.Create(ctx, recommendation)
}

// GetByID obtiene una recomendación por su ID
func (s *recommendationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recommendation, error) {
	return s.recommendationRepo.GetByID(ctx, id)
}

// GetByName obtiene una recomendación por su nombre
func (s *recommendationService) GetByName(ctx context.Context, name string) (*domain.Recommendation, error) {
	return s.recommendationRepo.GetByName(ctx, name)
}

// GetByUmbral obtiene recomendaciones por su umbral
func (s *recommendationService) GetByUmbral(ctx context.Context, umbral string) ([]*domain.Recommendation, error) {
	return s.recommendationRepo.GetByUmbral(ctx, umbral)
}

// GetAll obtiene todas las recomendaciones
func (s *recommendationService) GetAll(ctx context.Context) ([]*domain.Recommendation, error) {
	return s.recommendationRepo.GetAll(ctx)
}

// Update actualiza una recomendación existente
func (s *recommendationService) Update(ctx context.Context, recommendation *domain.Recommendation) error {
	if err := recommendation.Validate(); err != nil {
		return err
	}
	return s.recommendationRepo.Update(ctx, recommendation)
}

// Delete elimina una recomendación por su ID
func (s *recommendationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.recommendationRepo.Delete(ctx, id)
}