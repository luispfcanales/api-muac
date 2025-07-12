package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IMeasurementRepository define las operaciones para el repositorio de mediciones
type IMeasurementRepository interface {
	Create(ctx context.Context, measurement *domain.Measurement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Measurement, error)
	GetAll(ctx context.Context) ([]*domain.Measurement, error)
	Update(ctx context.Context, measurement *domain.Measurement) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Measurement, error)
	GetByTagID(ctx context.Context, tagID uuid.UUID) ([]*domain.Measurement, error)
	GetByRecommendationID(ctx context.Context, recommendationID uuid.UUID) ([]*domain.Measurement, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Measurement, error)
}

// IMeasurementService define las operaciones del servicio para mediciones (ACTUALIZADO)
type IMeasurementService interface {
	Create(ctx context.Context, measurement *domain.Measurement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Measurement, error)
	GetAll(ctx context.Context) ([]*domain.Measurement, error)
	Update(ctx context.Context, measurement *domain.Measurement) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Measurement, error)
	GetByTagID(ctx context.Context, tagID uuid.UUID) ([]*domain.Measurement, error)
	GetByRecommendationID(ctx context.Context, recommendationID uuid.UUID) ([]*domain.Measurement, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Measurement, error)
	AssignTag(ctx context.Context, measurementID, tagID uuid.UUID) error
	AssignRecommendation(ctx context.Context, measurementID, recommendationID uuid.UUID) error

	// ============= NUEVO MÉTODO PARA AUTO-ASIGNACIÓN =============
	CreateWithAutoAssignment(ctx context.Context, muacValue float64, description string, patientID, userID uuid.UUID) (*domain.Measurement, error)
}
