package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// roleService implementa la interfaz RoleService
type roleService struct {
	roleRepo ports.IRoleRepository
}

// NewRoleService crea una nueva instancia de RoleService
func NewRoleService(roleRepo ports.IRoleRepository) ports.IRoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

// CreateRole crea un nuevo rol
func (s *roleService) CreateRole(ctx context.Context, name, description string) (*domain.Role, error) {
	role := domain.NewRole(name, description)

	if err := role.Validate(); err != nil {
		return nil, err
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID obtiene un rol por su ID
func (s *roleService) GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	fmt.Printf("Buscando rol con ID: %s\n", id.String())
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("Error al buscar rol: %v\n", err)
	}
	return role, err
}

// GetAllRoles obtiene todos los roles
func (s *roleService) GetAllRoles(ctx context.Context) ([]*domain.Role, error) {
	return s.roleRepo.GetAll(ctx)
}

// UpdateRole actualiza un rol existente
func (s *roleService) UpdateRole(ctx context.Context, id uuid.UUID, name, description string) (*domain.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	role.Update(name, description)

	if err := role.Validate(); err != nil {
		return nil, err
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRole elimina un rol por su ID
func (s *roleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.roleRepo.Delete(ctx, id)
}
