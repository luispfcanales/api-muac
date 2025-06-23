package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tag representa la entidad de etiqueta en el dominio
type Tag struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"column:NAME;type:varchar(100);not null"`
	Description string    `json:"description" gorm:"column:DESCRIPTION;type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Tag) TableName() string {
	return "tags"
}

// NewTag crea una nueva instancia de Tag
func NewTag(name, description string) *Tag {
	return &Tag{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// Validate valida que la etiqueta tenga los campos requeridos
func (t *Tag) Validate() error {
	if t.Name == "" {
		return ErrEmptyTagName
	}
	return nil
}

// Update actualiza los campos de la etiqueta
func (t *Tag) Update(name, description string) {
	t.Name = name
	t.Description = description
	t.UpdatedAt = time.Now()
}
