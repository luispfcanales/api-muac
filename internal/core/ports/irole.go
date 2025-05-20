package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IRoleRepository define las operaciones que debe implementar un repositorio de roles
type IRoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	GetAll(ctx context.Context) ([]*domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IRoleService define las operaciones que debe implementar un servicio de roles
type IRoleService interface {
	CreateRole(ctx context.Context, name, description string) (*domain.Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)
	UpdateRole(ctx context.Context, id uuid.UUID, name, description string) (*domain.Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
}
