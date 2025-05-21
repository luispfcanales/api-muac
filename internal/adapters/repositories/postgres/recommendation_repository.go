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

// recommendationRepository implementa la interfaz IRecommendationRepository usando GORM
type recommendationRepository struct {
	db *gorm.DB
}

// NewRecommendationRepository crea una nueva instancia de RecommendationRepository
func NewRecommendationRepository(db *gorm.DB) ports.IRecommendationRepository {
	return &recommendationRepository{
		db: db,
	}
}

// Create inserta una nueva recomendación en la base de datos
func (r *recommendationRepository) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	result := r.db.WithContext(ctx).Create(recommendation)
	if result.Error != nil {
		return fmt.Errorf("error al crear recomendación: %w", result.Error)
	}
	return nil
}

// GetByID obtiene una recomendación por su ID
func (r *recommendationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recommendation, error) {
	var recommendation domain.Recommendation
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&recommendation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecommendationNotFound
		}
		return nil, fmt.Errorf("error al obtener recomendación: %w", result.Error)
	}
	return &recommendation, nil
}

// GetByName obtiene una recomendación por su nombre
func (r *recommendationRepository) GetByName(ctx context.Context, name string) (*domain.Recommendation, error) {
	var recommendation domain.Recommendation
	result := r.db.WithContext(ctx).Where("NAME = ?", name).First(&recommendation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecommendationNotFound
		}
		return nil, fmt.Errorf("error al obtener recomendación por nombre: %w", result.Error)
	}
	return &recommendation, nil
}

// GetByUmbral obtiene recomendaciones por su umbral
func (r *recommendationRepository) GetByUmbral(ctx context.Context, umbral string) ([]*domain.Recommendation, error) {
	var recommendations []*domain.Recommendation
	result := r.db.WithContext(ctx).Where("RECOMMENDATION_UMBRAL = ?", umbral).Find(&recommendations)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener recomendaciones por umbral: %w", result.Error)
	}
	return recommendations, nil
}

// GetAll obtiene todas las recomendaciones
func (r *recommendationRepository) GetAll(ctx context.Context) ([]*domain.Recommendation, error) {
	var recommendations []*domain.Recommendation
	result := r.db.WithContext(ctx).Find(&recommendations)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener recomendaciones: %w", result.Error)
	}
	return recommendations, nil
}

// Update actualiza una recomendación existente
func (r *recommendationRepository) Update(ctx context.Context, recommendation *domain.Recommendation) error {
	result := r.db.WithContext(ctx).Save(recommendation)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar recomendación: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRecommendationNotFound
	}
	return nil
}

// Delete elimina una recomendación por su ID
func (r *recommendationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Recommendation{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar recomendación: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRecommendationNotFound
	}
	return nil
}