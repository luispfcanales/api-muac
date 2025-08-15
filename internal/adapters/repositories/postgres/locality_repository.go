package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// localityRepository implementa la interfaz ILocalityRepository usando GORM
type localityRepository struct {
	db *gorm.DB
}

// NewLocalityRepository crea una nueva instancia de LocalityRepository
func NewLocalityRepository(db *gorm.DB) ports.ILocalityRepository {
	return &localityRepository{
		db: db,
	}
}

// Create inserta una nueva localidad en la base de datos
func (r *localityRepository) Create(ctx context.Context, locality *domain.Locality) error {
	result := r.db.WithContext(ctx).Create(locality)
	if result.Error != nil {
		return fmt.Errorf("error al crear localidad: %w", result.Error)
	}
	return nil
}

// GetByID obtiene una localidad por su ID
func (r *localityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Locality, error) {
	var locality domain.Locality
	result := r.db.WithContext(ctx).Where("ID = ?", id).First(&locality)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLocalityNotFound
		}
		return nil, fmt.Errorf("error al obtener localidad: %w", result.Error)
	}
	return &locality, nil
}

// GetByName obtiene una localidad por su nombre
func (r *localityRepository) GetByName(ctx context.Context, name string) (*domain.Locality, error) {
	var locality domain.Locality
	result := r.db.WithContext(ctx).Where("NAME = ?", name).First(&locality)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLocalityNotFound
		}
		return nil, fmt.Errorf("error al obtener localidad por nombre: %w", result.Error)
	}
	return &locality, nil
}

// GetAll obtiene todas las localidades
func (r *localityRepository) GetAll(ctx context.Context) ([]*domain.Locality, error) {
	var localities []*domain.Locality
	result := r.db.WithContext(ctx).Find(&localities)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener localidades: %w", result.Error)
	}
	return localities, nil
}

// Update actualiza una localidad existente
func (r *localityRepository) Update(ctx context.Context, locality *domain.Locality) error {
	result := r.db.WithContext(ctx).Save(locality)
	if result.Error != nil {
		return fmt.Errorf("error al actualizar localidad: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrLocalityNotFound
	}
	return nil
}

// Delete elimina una localidad por su ID
func (r *localityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Locality{}, "ID = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar localidad: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrLocalityNotFound
	}
	return nil
}

func (r *localityRepository) FindNearby(ctx context.Context, lat, lng float64, radiusKm float64) ([]domain.Locality, error) {
	var allLocalities []domain.Locality
	var nearbyLocalities []domain.Locality

	// Primero obtener todas las localidades
	if err := r.db.WithContext(ctx).Find(&allLocalities).Error; err != nil {
		return nil, fmt.Errorf("error fetching localities: %w", err)
	}

	// Filtrar en memoria
	for _, loc := range allLocalities {
		locLat, err := strconv.ParseFloat(loc.Latitude, 64)
		if err != nil {
			continue // O manejar el error como prefieras
		}

		locLng, err := strconv.ParseFloat(loc.Longitude, 64)
		if err != nil {
			continue // O manejar el error como prefieras
		}

		distance := haversine(lat, lng, locLat, locLng)
		if distance <= radiusKm {
			nearbyLocalities = append(nearbyLocalities, loc)
		}
	}

	// Ordenar por distancia (opcional)
	sort.Slice(nearbyLocalities, func(i, j int) bool {
		latI, _ := strconv.ParseFloat(nearbyLocalities[i].Latitude, 64)
		lngI, _ := strconv.ParseFloat(nearbyLocalities[i].Longitude, 64)
		latJ, _ := strconv.ParseFloat(nearbyLocalities[j].Latitude, 64)
		lngJ, _ := strconv.ParseFloat(nearbyLocalities[j].Longitude, 64)

		distI := haversine(lat, lng, latI, lngI)
		distJ := haversine(lat, lng, latJ, lngJ)
		return distI < distJ
	})

	return nearbyLocalities, nil
}

// Función Haversine implementada en Go
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radio de la Tierra en km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
