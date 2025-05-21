package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IPatientRepository define las operaciones para el repositorio de pacientes
type IPatientRepository interface {
	Create(ctx context.Context, patient *domain.Patient) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error)
	GetByDNI(ctx context.Context, dni string) (*domain.Patient, error)
	GetAll(ctx context.Context) ([]*domain.Patient, error)
	Update(ctx context.Context, patient *domain.Patient) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByFatherID(ctx context.Context, fatherID uuid.UUID) ([]*domain.Patient, error)
	GetMeasurements(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error)
}

// IPatientService define las operaciones del servicio para pacientes
type IPatientService interface {
	Create(ctx context.Context, patient *domain.Patient) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error)
	GetByDNI(ctx context.Context, dni string) (*domain.Patient, error)
	GetAll(ctx context.Context) ([]*domain.Patient, error)
	Update(ctx context.Context, patient *domain.Patient) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByFatherID(ctx context.Context, fatherID uuid.UUID) ([]*domain.Patient, error)
	GetMeasurements(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error)
	AddMeasurement(ctx context.Context, patientID uuid.UUID, measurement *domain.Measurement) error
}
