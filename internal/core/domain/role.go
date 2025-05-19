package domain

import (
	"github.com/google/uuid"
	"time"
)

// Role representa la entidad de rol en el dominio
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
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