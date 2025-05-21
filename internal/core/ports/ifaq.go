package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IFAQRepository define las operaciones para el repositorio de preguntas frecuentes
type IFAQRepository interface {
	Create(ctx context.Context, faq *domain.FAQ) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FAQ, error)
	GetAll(ctx context.Context) ([]*domain.FAQ, error)
	Update(ctx context.Context, faq *domain.FAQ) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IFAQService define las operaciones del servicio para preguntas frecuentes
type IFAQService interface {
	Create(ctx context.Context, faq *domain.FAQ) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FAQ, error)
	GetAll(ctx context.Context) ([]*domain.FAQ, error)
	Update(ctx context.Context, faq *domain.FAQ) error
	Delete(ctx context.Context, id uuid.UUID) error
}