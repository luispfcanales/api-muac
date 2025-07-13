package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// UserService implementa la lógica de negocio para usuarios
type userService struct {
	userRepo ports.IUserRepository
	roleRepo ports.IRoleRepository
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo ports.IUserRepository, roleRepo ports.IRoleRepository) ports.IUserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetByUsernameOrEmail obtiene un usuario por su nombre de usuario o email
func (s *userService) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domain.User, error) {
	return s.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
}

// Create crea un nuevo usuario
func (s *userService) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	// Verificar que el rol existe
	if user.RoleID != uuid.Nil {
		_, err := s.roleRepo.GetByID(ctx, user.RoleID)
		if err != nil {
			return err
		}
	}

	return s.userRepo.Create(ctx, user)
}

// GetByID obtiene un usuario por su ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail obtiene un usuario por su email
func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetAll obtiene todos los usuarios, opcionalmente filtrados por localidad
func (s *userService) GetAll(ctx context.Context, localityID *uuid.UUID) ([]*domain.User, error) {
	return s.userRepo.GetAll(ctx, localityID)
}

// GetApoderados obtiene todos los usuarios con rol de APODERADO, opcionalmente filtrados por localidad
func (s *userService) GetApoderados(ctx context.Context, localityID *uuid.UUID) ([]*domain.User, error) {
	return s.userRepo.GetByRole(ctx, "APODERADO", localityID)
}

// Update actualiza un usuario existente
func (s *userService) Update(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	// Verificar que el rol existe
	if user.RoleID != uuid.Nil {
		_, err := s.roleRepo.GetByID(ctx, user.RoleID)
		if err != nil {
			return err
		}
	}

	return s.userRepo.Update(ctx, user)
}

// Delete elimina un usuario por su ID
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// UpdatePassword actualiza la contraseña de un usuario
func (s *userService) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.UpdatePassword(passwordHash)
	return s.userRepo.Update(ctx, user)
}

// UpdateRole actualiza el rol de un usuario
func (s *userService) UpdateRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verificar que el rol existe
	if roleID != uuid.Nil {
		_, err := s.roleRepo.GetByID(ctx, roleID)
		if err != nil {
			return err
		}
	}

	user.UpdateRole(roleID)
	return s.userRepo.Update(ctx, user)
}
