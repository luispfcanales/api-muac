package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role representa la entidad de rol en el dominio
type Role struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Role) TableName() string {
	return "roles"
}

// NewRole crea una nueva instancia de Role
func NewRole(name, description string) *Role {
	return &Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// Validate valida que el rol tenga los campos requeridos
func (r *Role) Validate() error {
	if r.Name == "" {
		return ErrEmptyRoleName
	}
	return nil
}

// Update actualiza los campos del rol
func (r *Role) Update(name, description string) {
	r.Name = name
	r.Description = description
	r.UpdatedAt = time.Now()
}
