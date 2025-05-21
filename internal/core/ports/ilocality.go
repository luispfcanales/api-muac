package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// ILocalityRepository define las operaciones para el repositorio de localidades
type ILocalityRepository interface {
	Create(ctx context.Context, locality *domain.Locality) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Locality, error)
	GetAll(ctx context.Context) ([]*domain.Locality, error)
	Update(ctx context.Context, locality *domain.Locality) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*domain.Locality, error)
}

// ILocalityService define las operaciones del servicio para localidades
type ILocalityService interface {
	Create(ctx context.Context, locality *domain.Locality) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Locality, error)
	GetAll(ctx context.Context) ([]*domain.Locality, error)
	Update(ctx context.Context, locality *domain.Locality) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*domain.Locality, error)
}