package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// tagService implementa la lógica de negocio para etiquetas
type tagService struct {
	tagRepo ports.ITagRepository
}

// NewTagService crea una nueva instancia de TagService
func NewTagService(tagRepo ports.ITagRepository) ports.ITagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

// Create crea una nueva etiqueta
func (s *tagService) Create(ctx context.Context, tag *domain.Tag) error {
	if err := tag.Validate(); err != nil {
		return err
	}
	return s.tagRepo.Create(ctx, tag)
}

// GetByID obtiene una etiqueta por su ID
func (s *tagService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	return s.tagRepo.GetByID(ctx, id)
}

// GetByName obtiene una etiqueta por su nombre
func (s *tagService) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	return s.tagRepo.GetByName(ctx, name)
}

// GetAll obtiene todas las etiquetas
func (s *tagService) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	return s.tagRepo.GetAll(ctx)
}

// Update actualiza una etiqueta existente
func (s *tagService) Update(ctx context.Context, tag *domain.Tag) error {
	if err := tag.Validate(); err != nil {
		return err
	}
	return s.tagRepo.Update(ctx, tag)
}

// Delete elimina una etiqueta por su ID
func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.tagRepo.Delete(ctx, id)
}