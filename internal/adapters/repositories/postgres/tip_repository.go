package postgres

import (
	"context"
	"fmt"

	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// tagRepository implementa la interfaz ITagRepository usando GORM
type tipRepository struct {
	db *gorm.DB
}

// NewTagRepository crea una nueva instancia de TagRepository
func NewTipRepository(db *gorm.DB) ports.ITipRepository {
	return &tipRepository{
		db: db,
	}
}

// GetAll obtiene todas las recetas de consejos
func (r *tipRepository) GetAll(ctx context.Context) ([]*domain.Tip, error) {
	var tips []*domain.Tip
	result := r.db.WithContext(ctx).Find(&tips)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener todas las recetas de consejos: %w", result.Error)
	}
	return tips, nil
}
