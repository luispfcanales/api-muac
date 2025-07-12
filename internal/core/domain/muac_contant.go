// domain/constants.go
package domain

import "fmt"

// ============= CÓDIGOS MUAC OFICIALES =============
const (
	MuacCodeRed    = "MUAC-R1" // < 11.5 cm - Desnutrición aguda severa (SAM)
	MuacCodeYellow = "MUAC-Y1" // 11.5-12.4 cm - Desnutrición aguda moderada (MAM)
	MuacCodeGreen  = "MUAC-G1" // ≥ 12.5 cm - Estado nutricional adecuado
	MuacCodeFollow = "MUAC-S1" // Seguimiento post-intervención
)

// ============= COLORES OFICIALES =============
const (
	ColorRed    = "#dc3545" // Rojo intenso - Urgente
	ColorYellow = "#ffc107" // Amarillo fuerte - Atención
	ColorGreen  = "#28a745" // Verde brillante - Normal
	ColorBlue   = "#17a2b8" // Azul - Seguimiento
	ColorGray   = "#6c757d" // Gris - Por defecto
)

// ============= PRIORIDADES =============
const (
	PriorityNormal    = 1 // Verde - Seguimiento normal
	PriorityAttention = 2 // Amarillo - Requiere atención
	PriorityUrgent    = 3 // Rojo - Urgente

	// Para Tags (rango extendido)
	PriorityLow      = 1  // Baja
	PriorityMedium   = 3  // Media
	PriorityHigh     = 5  // Alta
	PriorityExtreme  = 8  // Extrema
	PriorityCritical = 10 // Crítica
)

// ============= UMBRALES MUAC OFICIALES =============
const (
	MuacThresholdSevere   = 11.5 // < 11.5 cm = SAM
	MuacThresholdModerate = 12.4 // 11.5-12.4 cm = MAM
	MuacThresholdNormal   = 12.5 // ≥ 12.5 cm = Normal
)

// ============= ERRORES COMUNES =============
var (
	// Errores de Tag
	// ErrEmptyTagName       = fmt.Errorf("el nombre del tag no puede estar vacío")
	ErrInvalidTagColor    = fmt.Errorf("código de color inválido")
	ErrInvalidMuacCode    = fmt.Errorf("código MUAC inválido")
	ErrInvalidTagPriority = fmt.Errorf("prioridad de tag inválida (debe estar entre 1 y 10)")
	// ErrTagNotFound        = fmt.Errorf("tag no encontrado")

	// Errores de Recommendation
	// ErrEmptyRecommendationName = fmt.Errorf("el nombre de la recomendación no puede estar vacío")
	ErrInvalidPriority  = fmt.Errorf("prioridad inválida: debe estar entre 1 y 3")
	ErrInvalidMuacRange = fmt.Errorf("rango MUAC inválido: mínimo no puede ser mayor que máximo")
	// ErrRecommendationNotFound  = fmt.Errorf("recomendación no encontrada")

	// Errores generales
	// ErrInvalidMuacValue = fmt.Errorf("valor MUAC inválido: debe ser mayor a 0")
	ErrDatabaseError = fmt.Errorf("error en la base de datos")
)

// ============= FUNCIONES HELPER GLOBALES =============

// ClassifyMuacValue clasifica un valor MUAC según estándares OMS
func ClassifyMuacValue(muacValue float64) (muacCode, colorCode string, priority int) {
	switch {
	case muacValue >= MuacThresholdNormal:
		return MuacCodeGreen, ColorGreen, PriorityNormal
	case muacValue >= MuacThresholdSevere:
		return MuacCodeYellow, ColorYellow, PriorityAttention
	default:
		return MuacCodeRed, ColorRed, PriorityUrgent
	}
}

// IsValidHexColor valida si es un código de color hexadecimal válido
func IsValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}

	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}

// IsValidMuacCode valida si es un código MUAC válido
func IsValidMuacCode(muacCode string) bool {
	validCodes := []string{MuacCodeRed, MuacCodeYellow, MuacCodeGreen, MuacCodeFollow}
	for _, code := range validCodes {
		if muacCode == code {
			return true
		}
	}
	return false
}

// IsValidMuacValue valida si un valor MUAC es válido
func IsValidMuacValue(value float64) bool {
	return value > 0 && value <= 50 // Límites razonables para MUAC
}

// GetMuacRiskLevel obtiene el nivel de riesgo textual
func GetMuacRiskLevel(muacValue float64) string {
	muacCode, _, _ := ClassifyMuacValue(muacValue)
	switch muacCode {
	case MuacCodeRed:
		return "Desnutrición Aguda Severa (SAM)"
	case MuacCodeYellow:
		return "Desnutrición Aguda Moderada (MAM)"
	case MuacCodeGreen:
		return "Estado Nutricional Adecuado"
	default:
		return "Sin Clasificar"
	}
}
