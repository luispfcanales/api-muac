package domain

import (
	"time"

	"github.com/google/uuid"
)

// Locality representa la entidad de localidad en el dominio
type Locality struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"column:NAME;type:varchar(100);not null"`
	Latitude    string    `json:"latitude" gorm:"column:LATITUDE;type:varchar(100)"`
	Longitude   string    `json:"longitude" gorm:"column:LONGITUDE;type:varchar(100)"`
	Description string    `json:"description" gorm:"column:DESCRIPTION;type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Locality) TableName() string {
	return "localities"
}

// NewLocality crea una nueva instancia de Locality
func NewLocality(name, latitude, longitude, description string) *Locality {
	return &Locality{
		ID:          uuid.New(),
		Name:        name,
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// Validate valida que la localidad tenga los campos requeridos
func (l *Locality) Validate() error {
	if l.Name == "" {
		return ErrEmptyLocalityName
	}
	return nil
}

// Update actualiza los campos de la localidad
func (l *Locality) Update(name, location, description string) {
	l.Name = name
	l.Description = description
	l.UpdatedAt = time.Now()
}
