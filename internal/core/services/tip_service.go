package services

import (
	"context"

	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// tagService implementa la lógica de negocio para etiquetas
type tipService struct {
	tipRepo ports.ITipRepository
}

// NewTagService crea una nueva instancia de TagService
func NewTipService(tipRepo ports.ITipRepository) ports.ITipService {
	return &tipService{
		tipRepo: tipRepo,
	}
}

// List obtiene todas las recomendaciones
func (s *tipService) List(ctx context.Context, muaccode string) ([]*domain.Tip, error) {

	if muaccode == domain.MuacCodeRed || muaccode == domain.MuacCodeFollow {
		return []*domain.Tip{}, nil
	}
	return s.tipRepo.GetAll(ctx)
}
