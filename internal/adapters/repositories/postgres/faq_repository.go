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

// faqRepository implementa la interfaz IFAQRepository usando GORM
type faqRepository struct {
	db *gorm.DB
}

// NewFAQRepository crea una nueva instancia de FAQRepository
func NewFAQRepository(db *gorm.DB) ports.IFAQRepository {
	return &faqRepository{
		db: db,
	}
}

// Create inserta una nueva FAQ en la base de datos
func (r *faqRepository) Create(ctx context.Context, faq *domain.FAQ) error {
	result := r.db.WithContext(ctx).Create(faq)
	if result.Error != nil {
		return fmt.Errorf("error al crear FAQ: %w", result.Error)
	}
	return nil
}

// GetByID obtiene una FAQ por su ID
func (r *faqRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.FAQ, error) {
	var faq domain.FAQ
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&faq)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFAQNotFound
		}
		return nil, fmt.Errorf("error al obtener FAQ: %w", result.Error)
	}
	return &faq, nil
}

// GetAll obtiene todas las FAQs
func (r *faqRepository) GetAll(ctx context.Context) ([]*domain.FAQ, error) {
	var faqs []*domain.FAQ
	result := r.db.WithContext(ctx).Find(&faqs)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener FAQs: %w", result.Error)
	}
	return faqs, nil
}

// Update actualiza una FAQ existente
func (r *faqRepository) Update(ctx context.Context, faq *domain.FAQ) error {
	result := r.db.WithContext(ctx).Save(faq)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar FAQ: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrFAQNotFound
	}
	return nil
}

// Delete elimina una FAQ por su ID
func (r *faqRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.FAQ{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar FAQ: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrFAQNotFound
	}
	return nil
}