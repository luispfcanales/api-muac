package domain

import "errors"

// Errores comunes del dominio
var (
	// Role errors
	ErrEmptyRoleName = errors.New("el nombre del rol no puede estar vacío")
	ErrRoleNotFound  = errors.New("rol no encontrado")

	// Locality errors
	ErrEmptyLocalityName     = errors.New("el nombre de la localidad no puede estar vacío")
	ErrEmptyLocalityLocation = errors.New("la ubicación de la localidad no puede estar vacía")
	ErrLocalityNotFound      = errors.New("localidad no encontrada")

	// Patient errors
	ErrEmptyPatientName        = errors.New("el nombre del paciente no puede estar vacío")
	ErrEmptyPatientLastName    = errors.New("el apellido del paciente no puede estar vacío")
	ErrPatientDNIAlreadyExists = errors.New("el DNI del paciente ya está registrado")
	ErrPatientNotFound         = errors.New("paciente no encontrado")

	// Tag errors
	ErrEmptyTagName = errors.New("el nombre de la etiqueta no puede estar vacío")
	ErrTagNotFound  = errors.New("etiqueta no encontrada")

	// User errors
	ErrEmptyUserName     = errors.New("el nombre del usuario no puede estar vacío")
	ErrEmptyUserLastName = errors.New("el apellido del usuario no puede estar vacío")
	ErrEmptyUsername     = errors.New("el nombre de usuario no puede estar vacío")
	ErrEmptyUserEmail    = errors.New("el email del usuario no puede estar vacío")
	ErrEmptyUserPassword = errors.New("la contraseña del usuario no puede estar vacía")
	ErrUserNotFound      = errors.New("usuario no encontrado")

	// Recommendation errors
	ErrEmptyRecommendationName = errors.New("el nombre de la recomendación no puede estar vacío")
	ErrRecommendationNotFound  = errors.New("recomendación no encontrada")

	// Measurement errors
	ErrInvalidMuacValue    = errors.New("el valor MUAC debe ser mayor que cero")
	ErrEmptyPatientID      = errors.New("el ID del paciente no puede estar vacío")
	ErrEmptyUserID         = errors.New("el ID del usuario no puede estar vacío")
	ErrMeasurementNotFound = errors.New("medición no encontrada")

	// Notification errors
	ErrEmptyNotificationTitle = errors.New("el título de la notificación no puede estar vacío")
	ErrNotificationNotFound   = errors.New("notificación no encontrada")

	// FAQ errors
	ErrEmptyFAQQuestion = errors.New("la pregunta no puede estar vacía")
	ErrEmptyFAQAnswer   = errors.New("la respuesta no puede estar vacía")
	ErrFAQNotFound      = errors.New("FAQ no encontrada")
)
