package services

import (
	"context"

	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// recipeService implementa la lógica de negocio para recetas
type recipeService struct {
	recipeRepo ports.IRecipeRepository
}

// NewRecipeService crea una nueva instancia de RecipeService
func NewRecipeService(recipeRepo ports.IRecipeRepository) ports.IRecipeService {
	return &recipeService{
		recipeRepo: recipeRepo,
	}
}

// ListRecipesByAge obtiene todas las recetas por edad
func (s *recipeService) ListRecipesByAge(ctx context.Context, age float32) ([]*domain.Recipe, error) {
	// Si la edad está fuera de los rangos válidos, retornar arreglo vacío
	if age < 0.5 || age > 5.0 {
		return []*domain.Recipe{}, nil
	}

	// Validar que la edad tenga formato válido (opcional, puedes quitarlo si no lo necesitas)
	if float32(int(age*10))/10 != age {
		return []*domain.Recipe{}, nil
	}

	return s.recipeRepo.GetRecipesByAge(ctx, age)
}
