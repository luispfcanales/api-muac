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

// tagRepository implementa la interfaz ITagRepository usando GORM
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository crea una nueva instancia de TagRepository
func NewTagRepository(db *gorm.DB) ports.ITagRepository {
	return &tagRepository{
		db: db,
	}
}

// Create inserta una nueva etiqueta en la base de datos
func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	result := r.db.WithContext(ctx).Create(tag)
	if result.Error != nil {
		return fmt.Errorf("error al crear etiqueta: %w", result.Error)
	}
	return nil
}

// GetByID obtiene una etiqueta por su ID
func (r *tagRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	var tag domain.Tag
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTagNotFound
		}
		return nil, fmt.Errorf("error al obtener etiqueta: %w", result.Error)
	}
	return &tag, nil
}

// GetByName obtiene una etiqueta por su nombre
func (r *tagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	var tag domain.Tag
	result := r.db.WithContext(ctx).Where("NAME = ?", name).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTagNotFound
		}
		return nil, fmt.Errorf("error al obtener etiqueta por nombre: %w", result.Error)
	}
	return &tag, nil
}

// GetAll obtiene todas las etiquetas
func (r *tagRepository) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	var tags []*domain.Tag
	result := r.db.WithContext(ctx).Find(&tags)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener etiquetas: %w", result.Error)
	}
	return tags, nil
}

// Update actualiza una etiqueta existente
func (r *tagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	result := r.db.WithContext(ctx).Save(tag)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar etiqueta: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrTagNotFound
	}
	return nil
}

// Delete elimina una etiqueta por su ID
func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Tag{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar etiqueta: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrTagNotFound
	}
	return nil
}