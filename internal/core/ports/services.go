package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// RoleService define las operaciones que debe implementar un servicio de roles
type RoleService interface {
	CreateRole(ctx context.Context, name, description string) (*domain.Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)
	UpdateRole(ctx context.Context, id uuid.UUID, name, description string) (*domain.Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
}
