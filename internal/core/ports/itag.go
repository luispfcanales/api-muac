package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// ITagRepository define las operaciones para el repositorio de etiquetas
type ITagRepository interface {
	Create(ctx context.Context, tag *domain.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	GetAll(ctx context.Context) ([]*domain.Tag, error)
	Update(ctx context.Context, tag *domain.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*domain.Tag, error)
}

// ITagService define las operaciones del servicio para etiquetas
type ITagService interface {
	Create(ctx context.Context, tag *domain.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	GetAll(ctx context.Context) ([]*domain.Tag, error)
	Update(ctx context.Context, tag *domain.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*domain.Tag, error)
}