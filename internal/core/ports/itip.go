package ports

import (
	"context"

	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// ITipRepository define las operaciones para el repositorio de etiquetas
type ITipRepository interface {
	GetAll(ctx context.Context) ([]*domain.Tip, error)
}

// ITipService define las operaciones del servicio para etiquetas
type ITipService interface {
	List(ctx context.Context, muaccode string) ([]*domain.Tip, error)
}
