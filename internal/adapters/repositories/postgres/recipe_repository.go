package postgres

import (
	"context"

	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// recipeRepository implementa la interfaz IRecipeRepository usando GORM
type recipeRepository struct {
	db *gorm.DB
}

// NewRecipeRepository crea una nueva instancia de RecipeRepository
func NewRecipeRepository(db *gorm.DB) ports.IRecipeRepository {
	return &recipeRepository{
		db: db,
	}
}

// GetRecipesByAge obtiene todas las recetas por edad
func (r *recipeRepository) GetRecipesByAge(ctx context.Context, age float64) ([]*domain.Recipe, error) {
	var recipes []*domain.Recipe
	err := r.db.WithContext(ctx).
		Where("min_age_years <= ? AND max_age_years > ?", age, age).
		Find(&recipes).Error
	if err != nil {
		return nil, err
	}
	return recipes, nil
}
