package ports

import (
	"context"

	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IRecipeRepository define las operaciones para el repositorio de recetas
type IRecipeRepository interface {
	GetRecipesByAge(ctx context.Context, age float32) ([]*domain.Recipe, error)
}

// IRecipeService define las operaciones del servicio para recetas
type IRecipeService interface {
	ListRecipesByAge(ctx context.Context, age float32) ([]*domain.Recipe, error)
}
