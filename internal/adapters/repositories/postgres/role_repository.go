package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"gorm.io/gorm"
)

// roleRepository implementa la interfaz RoleRepository usando GORM
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository crea una nueva instancia de RoleRepository
func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

// Create inserta un nuevo rol en la base de datos
func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	result := r.db.WithContext(ctx).Create(role)
	if result.Error != nil {
		return fmt.Errorf("error al crear rol: %w", result.Error)
	}
	return nil
}

// GetByID obtiene un rol por su ID
func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, fmt.Errorf("error al obtener rol: %w", result.Error)
	}
	return &role, nil
}

// GetAll obtiene todos los roles
func (r *roleRepository) GetAll(ctx context.Context) ([]*domain.Role, error) {
	var roles []*domain.Role
	result := r.db.WithContext(ctx).Find(&roles)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener roles: %w", result.Error)
	}
	return roles, nil
}

// Update actualiza un rol existente
func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	result := r.db.WithContext(ctx).Save(role)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar rol: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRoleNotFound
	}
	return nil
}

// Delete elimina un rol por su ID
func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Role{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar rol: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRoleNotFound
	}
	return nil
}
