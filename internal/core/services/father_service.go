package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// fatherService implementa la lógica de negocio para padres
type fatherService struct {
	fatherRepo   ports.IFatherRepository
	roleRepo     ports.IRoleRepository
	localityRepo ports.ILocalityRepository
	patientRepo  ports.IPatientRepository
}

// NewFatherService crea una nueva instancia de FatherService
func NewFatherService(
	fatherRepo ports.IFatherRepository,
	roleRepo ports.IRoleRepository,
	localityRepo ports.ILocalityRepository,
	patientRepo ports.IPatientRepository,
) ports.IFatherService {
	return &fatherService{
		fatherRepo:   fatherRepo,
		roleRepo:     roleRepo,
		localityRepo: localityRepo,
		patientRepo:  patientRepo,
	}
}

// Create crea un nuevo padre
func (s *fatherService) Create(ctx context.Context, father *domain.Father) error {
	if err := father.Validate(); err != nil {
		return err
	}

	// Verificar que el rol existe
	if father.RoleID != uuid.Nil {
		_, err := s.roleRepo.GetByID(ctx, father.RoleID)
		if err != nil {
			return err
		}
	}

	// Verificar que la localidad existe
	if father.LocalityID != uuid.Nil {
		_, err := s.localityRepo.GetByID(ctx, father.LocalityID)
		if err != nil {
			return err
		}
	}

	// Verificar que el paciente existe
	if father.PatientID != uuid.Nil {
		_, err := s.patientRepo.GetByID(ctx, father.PatientID)
		if err != nil {
			return err
		}
	}

	return s.fatherRepo.Create(ctx, father)
}

// GetByID obtiene un padre por su ID
func (s *fatherService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Father, error) {
	return s.fatherRepo.GetByID(ctx, id)
}

// GetByEmail obtiene un padre por su email
func (s *fatherService) GetByEmail(ctx context.Context, email string) (*domain.Father, error) {
	return s.fatherRepo.GetByEmail(ctx, email)
}

// GetByDNI obtiene un padre por su DNI
func (s *fatherService) GetByDNI(ctx context.Context, dni int) (*domain.Father, error) {
	return s.fatherRepo.GetByDNI(ctx, dni)
}

// GetByPatientID obtiene padres por ID de paciente
func (s *fatherService) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Father, error) {
	return s.fatherRepo.GetByPatientID(ctx, patientID)
}

// GetByLocalityID obtiene padres por ID de localidad
func (s *fatherService) GetByLocalityID(ctx context.Context, localityID uuid.UUID) ([]*domain.Father, error) {
	return s.fatherRepo.GetByLocalityID(ctx, localityID)
}

// GetAll obtiene todos los padres
func (s *fatherService) GetAll(ctx context.Context) ([]*domain.Father, error) {
	return s.fatherRepo.GetAll(ctx)
}

// Update actualiza un padre existente
func (s *fatherService) Update(ctx context.Context, father *domain.Father) error {
	if err := father.Validate(); err != nil {
		return err
	}

	// Verificar que el rol existe
	if father.RoleID != uuid.Nil {
		_, err := s.roleRepo.GetByID(ctx, father.RoleID)
		if err != nil {
			return err
		}
	}

	// Verificar que la localidad existe
	if father.LocalityID != uuid.Nil {
		_, err := s.localityRepo.GetByID(ctx, father.LocalityID)
		if err != nil {
			return err
		}
	}

	// Verificar que el paciente existe
	if father.PatientID != uuid.Nil {
		_, err := s.patientRepo.GetByID(ctx, father.PatientID)
		if err != nil {
			return err
		}
	}

	return s.fatherRepo.Update(ctx, father)
}

// Delete elimina un padre por su ID
func (s *fatherService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.fatherRepo.Delete(ctx, id)
}

// UpdatePassword actualiza la contraseña de un padre
func (s *fatherService) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	father, err := s.fatherRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	father.UpdatePassword(passwordHash)
	return s.fatherRepo.Update(ctx, father)
}

// UpdateActive actualiza el estado activo de un padre
func (s *fatherService) UpdateActive(ctx context.Context, id uuid.UUID, active bool) error {
	father, err := s.fatherRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	father.UpdateActive(active)
	return s.fatherRepo.Update(ctx, father)
}
