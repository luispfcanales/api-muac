package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IFatherRepository define las operaciones para el repositorio de padres
type IFatherRepository interface {
	Create(ctx context.Context, father *domain.Father) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Father, error)
	GetAll(ctx context.Context) ([]*domain.Father, error)
	Update(ctx context.Context, father *domain.Father) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*domain.Father, error)
	GetByDNI(ctx context.Context, dni int) (*domain.Father, error)
	GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Father, error)
	GetByLocalityID(ctx context.Context, localityID uuid.UUID) ([]*domain.Father, error)
}

// IFatherService define las operaciones del servicio para padres
type IFatherService interface {
	Create(ctx context.Context, father *domain.Father) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Father, error)
	GetAll(ctx context.Context) ([]*domain.Father, error)
	Update(ctx context.Context, father *domain.Father) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*domain.Father, error)
	GetByDNI(ctx context.Context, dni int) (*domain.Father, error)
	GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Father, error)
	GetByLocalityID(ctx context.Context, localityID uuid.UUID) ([]*domain.Father, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateActive(ctx context.Context, id uuid.UUID, active bool) error
}