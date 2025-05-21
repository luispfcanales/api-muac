package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IRecommendationRepository define las operaciones para el repositorio de recomendaciones
type IRecommendationRepository interface {
	Create(ctx context.Context, recommendation *domain.Recommendation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Recommendation, error)
	GetAll(ctx context.Context) ([]*domain.Recommendation, error)
	Update(ctx context.Context, recommendation *domain.Recommendation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*domain.Recommendation, error)
	GetByUmbral(ctx context.Context, umbral string) ([]*domain.Recommendation, error)
}

// IRecommendationService define las operaciones del servicio para recomendaciones
type IRecommendationService interface {
	Create(ctx context.Context, recommendation *domain.Recommendation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Recommendation, error)
	GetAll(ctx context.Context) ([]*domain.Recommendation, error)
	Update(ctx context.Context, recommendation *domain.Recommendation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*domain.Recommendation, error)
	GetByUmbral(ctx context.Context, umbral string) ([]*domain.Recommendation, error)
}