package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// userRepository implementa la interfaz UserRepository usando GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *gorm.DB) ports.IUserRepository {
	return &userRepository{
		db: db,
	}
}

// GetByUsername obtiene un usuario por su nombre de usuario
func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Locality").
		Preload("Patients").
		Where(`username = ? OR email = ?`, usernameOrEmail, usernameOrEmail).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", result.Error)
	}
	return &user, nil
}

// Create inserta un nuevo usuario en la base de datos
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("error al crear usuario: %w", result.Error)
	}
	return nil
}

// GetByID obtiene un usuario por su ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Locality").
		Preload("Patients").
		Preload("Patients.Measurements").
		Where("ID = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", result.Error)
	}
	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error al obtener usuario por email: %w", result.Error)
	}
	return &user, nil
}

// GetAll obtiene todos los usuarios con sus relaciones, opcionalmente filtrados por localidad
func (r *userRepository) GetAll(ctx context.Context, localityID *uuid.UUID) ([]*domain.User, error) {
	var users []*domain.User

	// Corregir los preloads - no existe "Recommendations" (plural) en Measurement
	query := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Locality").
		Preload("Patients").
		Preload("Patients.Measurements", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC") // Ordenar las mediciones por paciente
		}).
		Preload("Patients.Measurements.Tag").           // Singular, no Tags
		Preload("Patients.Measurements.Recommendation") // Singular, no Recommendations

	// Aplicar filtro por localidad si se proporciona
	if localityID != nil {
		query = query.Where("locality_id = ?", *localityID)
	}

	result := query.Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", result.Error)
	}
	return users, nil
}

// Update actualiza un usuario existente
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar usuario: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// Delete elimina un usuario por su ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar usuario: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
