package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// localityService implementa la lógica de negocio para localidades
type localityService struct {
	localityRepo ports.ILocalityRepository
}

// NewLocalityService crea una nueva instancia de LocalityService
func NewLocalityService(localityRepo ports.ILocalityRepository) ports.ILocalityService {
	return &localityService{
		localityRepo: localityRepo,
	}
}

// Create crea una nueva localidad
func (s *localityService) Create(ctx context.Context, locality *domain.Locality) error {
	if err := locality.Validate(); err != nil {
		return err
	}
	return s.localityRepo.Create(ctx, locality)
}

// GetByID obtiene una localidad por su ID
func (s *localityService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Locality, error) {
	return s.localityRepo.GetByID(ctx, id)
}

// GetByName obtiene una localidad por su nombre
func (s *localityService) GetByName(ctx context.Context, name string) (*domain.Locality, error) {
	return s.localityRepo.GetByName(ctx, name)
}

// GetAll obtiene todas las localidades
func (s *localityService) GetAll(ctx context.Context) ([]*domain.Locality, error) {
	return s.localityRepo.GetAll(ctx)
}

// Update actualiza una localidad existente
func (s *localityService) Update(ctx context.Context, locality *domain.Locality) error {
	if err := locality.Validate(); err != nil {
		return err
	}
	return s.localityRepo.Update(ctx, locality)
}

// Delete elimina una localidad por su ID
func (s *localityService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.localityRepo.Delete(ctx, id)
}