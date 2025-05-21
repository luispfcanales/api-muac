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

// fatherRepository implementa la interfaz IFatherRepository usando GORM
type fatherRepository struct {
	db *gorm.DB
}

// NewFatherRepository crea una nueva instancia de FatherRepository
func NewFatherRepository(db *gorm.DB) ports.IFatherRepository {
	return &fatherRepository{
		db: db,
	}
}

// Create inserta un nuevo padre en la base de datos
func (r *fatherRepository) Create(ctx context.Context, father *domain.Father) error {
	result := r.db.WithContext(ctx).Create(father)
	if result.Error != nil {
		return fmt.Errorf("error al crear padre: %w", result.Error)
	}
	return nil
}

// GetByID obtiene un padre por su ID
func (r *fatherRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Father, error) {
	var father domain.Father
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&father)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFatherNotFound
		}
		return nil, fmt.Errorf("error al obtener padre: %w", result.Error)
	}
	return &father, nil
}

// GetByEmail obtiene un padre por su email
func (r *fatherRepository) GetByEmail(ctx context.Context, email string) (*domain.Father, error) {
	var father domain.Father
	result := r.db.WithContext(ctx).Where("EMAIL = ?", email).First(&father)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFatherNotFound
		}
		return nil, fmt.Errorf("error al obtener padre por email: %w", result.Error)
	}
	return &father, nil
}

// GetByDNI obtiene un padre por su DNI
func (r *fatherRepository) GetByDNI(ctx context.Context, dni int) (*domain.Father, error) {
	var father domain.Father
	result := r.db.WithContext(ctx).Where("DNI = ?", dni).First(&father)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFatherNotFound
		}
		return nil, fmt.Errorf("error al obtener padre por DNI: %w", result.Error)
	}
	return &father, nil
}

// GetByPatientID obtiene padres por ID de paciente
func (r *fatherRepository) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Father, error) {
	var fathers []*domain.Father
	result := r.db.WithContext(ctx).Where("PATIENT_ID = ?", patientID).Find(&fathers)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener padres por ID de paciente: %w", result.Error)
	}
	return fathers, nil
}

// GetByLocalityID obtiene padres por ID de localidad
func (r *fatherRepository) GetByLocalityID(ctx context.Context, localityID uuid.UUID) ([]*domain.Father, error) {
	var fathers []*domain.Father
	result := r.db.WithContext(ctx).Where("LOCALITY_ID = ?", localityID).Find(&fathers)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener padres por ID de localidad: %w", result.Error)
	}
	return fathers, nil
}

// GetAll obtiene todos los padres
func (r *fatherRepository) GetAll(ctx context.Context) ([]*domain.Father, error) {
	var fathers []*domain.Father
	result := r.db.WithContext(ctx).Find(&fathers)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener padres: %w", result.Error)
	}
	return fathers, nil
}

// Update actualiza un padre existente
func (r *fatherRepository) Update(ctx context.Context, father *domain.Father) error {
	result := r.db.WithContext(ctx).Save(father)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar padre: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrFatherNotFound
	}
	return nil
}

// Delete elimina un padre por su ID
func (r *fatherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Father{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar padre: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrFatherNotFound
	}
	return nil
}