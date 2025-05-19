package domain

import "errors"

// Errores del dominio
var (
	ErrRoleNotFound   = errors.New("rol no encontrado")
	ErrEmptyRoleName  = errors.New("el nombre del rol no puede estar vacío")
	ErrRoleIDRequired = errors.New("se requiere el ID del rol")
	ErrInternalServer = errors.New("error interno del servidor")
)