package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// faqService implementa la lógica de negocio para preguntas frecuentes
type faqService struct {
	faqRepo ports.IFAQRepository
}

// NewFAQService crea una nueva instancia de FAQService
func NewFAQService(faqRepo ports.IFAQRepository) ports.IFAQService {
	return &faqService{
		faqRepo: faqRepo,
	}
}

// Create crea una nueva FAQ
func (s *faqService) Create(ctx context.Context, faq *domain.FAQ) error {
	if err := faq.Validate(); err != nil {
		return err
	}
	return s.faqRepo.Create(ctx, faq)
}

// GetByID obtiene una FAQ por su ID
func (s *faqService) GetByID(ctx context.Context, id uuid.UUID) (*domain.FAQ, error) {
	return s.faqRepo.GetByID(ctx, id)
}

// GetAll obtiene todas las FAQs
func (s *faqService) GetAll(ctx context.Context) ([]*domain.FAQ, error) {
	return s.faqRepo.GetAll(ctx)
}

// Update actualiza una FAQ existente
func (s *faqService) Update(ctx context.Context, faq *domain.FAQ) error {
	if err := faq.Validate(); err != nil {
		return err
	}
	return s.faqRepo.Update(ctx, faq)
}

// Delete elimina una FAQ por su ID
func (s *faqService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.faqRepo.Delete(ctx, id)
}