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

// localityRepository implementa la interfaz ILocalityRepository usando GORM
type localityRepository struct {
	db *gorm.DB
}

// NewLocalityRepository crea una nueva instancia de LocalityRepository
func NewLocalityRepository(db *gorm.DB) ports.ILocalityRepository {
	return &localityRepository{
		db: db,
	}
}

// Create inserta una nueva localidad en la base de datos
func (r *localityRepository) Create(ctx context.Context, locality *domain.Locality) error {
	result := r.db.WithContext(ctx).Create(locality)
	if result.Error != nil {
		return fmt.Errorf("error al crear localidad: %w", result.Error)
	}
	return nil
}

// GetByID obtiene una localidad por su ID
func (r *localityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Locality, error) {
	var locality domain.Locality
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&locality)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLocalityNotFound
		}
		return nil, fmt.Errorf("error al obtener localidad: %w", result.Error)
	}
	return &locality, nil
}

// GetByName obtiene una localidad por su nombre
func (r *localityRepository) GetByName(ctx context.Context, name string) (*domain.Locality, error) {
	var locality domain.Locality
	result := r.db.WithContext(ctx).Where("NAME = ?", name).First(&locality)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLocalityNotFound
		}
		return nil, fmt.Errorf("error al obtener localidad por nombre: %w", result.Error)
	}
	return &locality, nil
}

// GetAll obtiene todas las localidades
func (r *localityRepository) GetAll(ctx context.Context) ([]*domain.Locality, error) {
	var localities []*domain.Locality
	result := r.db.WithContext(ctx).Find(&localities)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener localidades: %w", result.Error)
	}
	return localities, nil
}

// Update actualiza una localidad existente
func (r *localityRepository) Update(ctx context.Context, locality *domain.Locality) error {
	result := r.db.WithContext(ctx).Save(locality)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar localidad: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrLocalityNotFound
	}
	return nil
}

// Delete elimina una localidad por su ID
func (r *localityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Locality{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar localidad: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrLocalityNotFound
	}
	return nil
}