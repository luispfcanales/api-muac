package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IUserRepository define las operaciones para el repositorio de usuarios
type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domain.User, error)
	GetAll(ctx context.Context, localityID *uuid.UUID) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByRole(ctx context.Context, roleName string, localityID *uuid.UUID) ([]*domain.User, error)
}

// IUserService define las operaciones del servicio para usuarios
type IUserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domain.User, error)
	GetAll(ctx context.Context, localityID *uuid.UUID) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) error
	GetApoderados(ctx context.Context, localityID *uuid.UUID) ([]*domain.User, error)
}
