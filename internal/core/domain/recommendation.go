// domain/recommendation.go
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Recommendation representa la entidad de recomendación nutricional MUAC en el dominio
type Recommendation struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name                 string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Description          string    `json:"description" gorm:"column:description;type:text;not null"`
	RecommendationUmbral string    `json:"recommendation_umbral" gorm:"column:recommendation_umbral;type:varchar(255)"`

	// Campos MUAC específicos
	MinValue  *float64 `json:"min_value,omitempty" gorm:"column:min_value;type:decimal(10,2)"`
	MaxValue  *float64 `json:"max_value,omitempty" gorm:"column:max_value;type:decimal(10,2)"`
	Priority  int      `json:"priority" gorm:"column:priority;type:int;default:1"`
	Active    bool     `json:"active" gorm:"column:active;default:true"`
	ColorCode string   `json:"color_code" gorm:"column:color_code;type:varchar(20)"`
	MuacCode  string   `json:"muac_code" gorm:"column:muac_code;type:varchar(10)"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Recommendation) TableName() string {
	return "recommendations"
}

// ============= CONSTRUCTORES =============

// NewRecommendation crea una nueva recomendación básica
func NewRecommendation(name, description, umbral string) *Recommendation {
	return &Recommendation{
		ID:                   uuid.New(),
		Name:                 name,
		Description:          description,
		RecommendationUmbral: umbral,
		Priority:             PriorityNormal,
		Active:               true,
		ColorCode:            ColorGray,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// NewMuacRecommendation crea una recomendación específica para MUAC
func NewMuacRecommendation(name, description string, minValue, maxValue *float64, priority int, colorCode, muacCode string) *Recommendation {
	rec := &Recommendation{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		MinValue:    minValue,
		MaxValue:    maxValue,
		Priority:    priority,
		Active:      true,
		ColorCode:   colorCode,
		MuacCode:    muacCode,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Generar umbral automáticamente
	rec.RecommendationUmbral = rec.generateUmbralText()

	// Validar y usar color por defecto si es necesario
	if !IsValidHexColor(colorCode) {
		rec.ColorCode = rec.getColorForMuacCode()
	}

	return rec
}

// ============= VALIDACIÓN =============

// Validate valida que la recomendación tenga los campos requeridos
func (r *Recommendation) Validate() error {
	if r.Name == "" {
		return ErrEmptyRecommendationName
	}

	if r.Description == "" {
		return fmt.Errorf("la descripción de la recomendación es requerida")
	}

	if r.MinValue != nil && r.MaxValue != nil && *r.MinValue > *r.MaxValue {
		return ErrInvalidMuacRange
	}

	if r.MinValue != nil && !IsValidMuacValue(*r.MinValue) {
		return fmt.Errorf("%w: valor mínimo inválido", ErrInvalidMuacValue)
	}

	if r.MaxValue != nil && !IsValidMuacValue(*r.MaxValue) {
		return fmt.Errorf("%w: valor máximo inválido", ErrInvalidMuacValue)
	}

	if r.Priority < 1 || r.Priority > 3 {
		return ErrInvalidPriority
	}

	if r.MuacCode != "" && !IsValidMuacCode(r.MuacCode) {
		return fmt.Errorf("%w: %s", ErrInvalidMuacCode, r.MuacCode)
	}

	if r.ColorCode != "" && !IsValidHexColor(r.ColorCode) {
		return fmt.Errorf("%w: %s", ErrInvalidTagColor, r.ColorCode)
	}

	return nil
}

// ============= MÉTODOS DE ACTUALIZACIÓN =============

// Update actualiza los campos básicos de la recomendación
func (r *Recommendation) Update(name, description, umbral string) {
	if name != "" {
		r.Name = name
	}
	if description != "" {
		r.Description = description
	}
	if umbral != "" {
		r.RecommendationUmbral = umbral
	}
	r.UpdatedAt = time.Now()
}

// UpdateMuacRecommendation actualiza una recomendación MUAC completa
func (r *Recommendation) UpdateMuacRecommendation(name, description string, minValue, maxValue *float64, priority int, colorCode, muacCode string) error {
	// Validar rangos antes de actualizar
	if minValue != nil && maxValue != nil && *minValue > *maxValue {
		return ErrInvalidMuacRange
	}

	if priority < 1 || priority > 3 {
		return ErrInvalidPriority
	}

	if muacCode != "" && !IsValidMuacCode(muacCode) {
		return fmt.Errorf("%w: %s", ErrInvalidMuacCode, muacCode)
	}

	if colorCode != "" && !IsValidHexColor(colorCode) {
		return fmt.Errorf("%w: %s", ErrInvalidTagColor, colorCode)
	}

	// Actualizar campos
	if name != "" {
		r.Name = name
	}
	if description != "" {
		r.Description = description
	}
	r.MinValue = minValue
	r.MaxValue = maxValue
	r.Priority = priority
	r.ColorCode = colorCode
	r.MuacCode = muacCode
	r.RecommendationUmbral = r.generateUmbralText()
	r.UpdatedAt = time.Now()

	return nil
}

// SetPriority establece la prioridad de la recomendación
func (r *Recommendation) SetPriority(priority int) error {
	if priority < 1 || priority > 3 {
		return ErrInvalidPriority
	}
	r.Priority = priority
	r.UpdatedAt = time.Now()
	return nil
}

// SetMuacRange establece el rango MUAC para la recomendación
func (r *Recommendation) SetMuacRange(minValue, maxValue *float64) error {
	if minValue != nil && maxValue != nil && *minValue > *maxValue {
		return ErrInvalidMuacRange
	}

	r.MinValue = minValue
	r.MaxValue = maxValue
	r.RecommendationUmbral = r.generateUmbralText()
	r.UpdatedAt = time.Now()
	return nil
}

// Activate activa la recomendación
func (r *Recommendation) Activate() {
	r.Active = true
	r.UpdatedAt = time.Now()
}

// Deactivate desactiva la recomendación
func (r *Recommendation) Deactivate() {
	r.Active = false
	r.UpdatedAt = time.Now()
}

// ============= MÉTODOS DE CONSULTA =============

// IsApplicableForMuac verifica si la recomendación aplica para un valor MUAC
func (r *Recommendation) IsApplicableForMuac(muacValue float64) bool {
	if !r.Active {
		return false
	}

	if !IsValidMuacValue(muacValue) {
		return false
	}

	if r.MinValue == nil && r.MaxValue == nil {
		return true
	}

	if r.MinValue != nil && muacValue < *r.MinValue {
		return false
	}

	if r.MaxValue != nil && muacValue >= *r.MaxValue {
		return false
	}

	return true
}

// IsUrgent verifica si la recomendación es urgente
func (r *Recommendation) IsUrgent() bool {
	return r.Priority >= PriorityUrgent
}

// IsNormal verifica si la recomendación es de prioridad normal
func (r *Recommendation) IsNormal() bool {
	return r.Priority == PriorityNormal
}

// HasMuacRange verifica si tiene rangos MUAC definidos
func (r *Recommendation) HasMuacRange() bool {
	return r.MinValue != nil || r.MaxValue != nil
}

// GetPriorityText retorna el texto de prioridad
func (r *Recommendation) GetPriorityText() string {
	switch r.Priority {
	case PriorityNormal:
		return "Normal"
	case PriorityAttention:
		return "Atención"
	case PriorityUrgent:
		return "Urgente"
	default:
		return "Desconocida"
	}
}

// GetUmbralDisplay retorna el umbral para mostrar en UI
func (r *Recommendation) GetUmbralDisplay() string {
	if r.RecommendationUmbral != "" {
		return r.RecommendationUmbral
	}
	return r.generateUmbralText()
}

// GetColorOrDefault retorna el color o un color por defecto
func (r *Recommendation) GetColorOrDefault() string {
	if r.ColorCode != "" {
		return r.ColorCode
	}
	return r.getColorForMuacCode()
}

// ============= FUNCIONES DE UTILIDAD PRIVADAS =============

// generateUmbralText genera texto del umbral basado en valores min/max
func (r *Recommendation) generateUmbralText() string {
	if r.MinValue == nil && r.MaxValue == nil {
		return "Todas las mediciones"
	}

	if r.MinValue == nil && r.MaxValue != nil {
		return fmt.Sprintf("< %.1f cm", *r.MaxValue)
	}

	if r.MinValue != nil && r.MaxValue == nil {
		return fmt.Sprintf("≥ %.1f cm", *r.MinValue)
	}

	return fmt.Sprintf("%.1f - %.1f cm", *r.MinValue, *r.MaxValue)
}

// getColorForMuacCode obtiene el color según el código MUAC
func (r *Recommendation) getColorForMuacCode() string {
	switch r.MuacCode {
	case MuacCodeRed:
		return ColorRed
	case MuacCodeYellow:
		return ColorYellow
	case MuacCodeGreen:
		return ColorGreen
	case MuacCodeFollow:
		return ColorBlue
	default:
		return ColorGray
	}
}

// ============= FUNCIONES HELPER GLOBALES =============

// GetRecommendationForMuacValue obtiene la recomendación apropiada para un valor MUAC
func GetRecommendationForMuacValue(muacValue float64, recommendations []*Recommendation) *Recommendation {
	// Ordenar por prioridad (urgente primero)
	for priority := PriorityUrgent; priority >= PriorityNormal; priority-- {
		for _, rec := range recommendations {
			if rec.Priority == priority && rec.IsApplicableForMuac(muacValue) {
				return rec
			}
		}
	}
	return nil
}

// FilterActiveRecommendations filtra solo recomendaciones activas
func FilterActiveRecommendations(recommendations []*Recommendation) []*Recommendation {
	var activeRecs []*Recommendation
	for _, rec := range recommendations {
		if rec.Active {
			activeRecs = append(activeRecs, rec)
		}
	}
	return activeRecs
}

// SortRecommendationsByPriority ordena recomendaciones por prioridad (mayor a menor)
func SortRecommendationsByPriority(recommendations []*Recommendation) []*Recommendation {
	sorted := make([]*Recommendation, len(recommendations))
	copy(sorted, recommendations)

	// Bubble sort simple - para datasets pequeños es suficiente
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Priority < sorted[j+1].Priority {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// GetRecommendationsByMuacCode obtiene recomendaciones por código MUAC
func GetRecommendationsByMuacCode(muacCode string, recommendations []*Recommendation) []*Recommendation {
	var filtered []*Recommendation
	for _, rec := range recommendations {
		if rec.MuacCode == muacCode && rec.Active {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

// CreateDefaultMuacRecommendations crea las recomendaciones por defecto del sistema
func CreateDefaultMuacRecommendations() []*Recommendation {
	valorSevere := MuacThresholdSevere
	valorModerate := MuacThresholdModerate
	valorNormal := MuacThresholdNormal

	return []*Recommendation{
		NewMuacRecommendation(
			"🚨 ALERTA ROJA - Acción Urgente Requerida",
			"⚠️ Esta medición indica DESNUTRICIÓN AGUDA SEVERA (SAM). Requiere atención médica URGENTE.",
			nil, &valorSevere,
			PriorityUrgent,
			ColorRed,
			MuacCodeRed,
		),
		NewMuacRecommendation(
			"🟡 ALERTA AMARILLA - Zona de Riesgo Nutricional",
			"🟡 El niño/a está en RIESGO NUTRICIONAL (MAM). Requiere mejoras en alimentación.",
			&valorSevere, &valorModerate,
			PriorityAttention,
			ColorYellow,
			MuacCodeYellow,
		),
		NewMuacRecommendation(
			"✅ ZONA VERDE - Estado Nutricional Adecuado",
			"✅ ¡Excelente! El niño/a tiene BUEN ESTADO NUTRICIONAL. Mantener cuidados.",
			&valorNormal, nil,
			PriorityNormal,
			ColorGreen,
			MuacCodeGreen,
		),
		NewMuacRecommendation(
			"📋 Seguimiento Post-Intervención",
			"📋 Paciente en proceso de RECUPERACIÓN NUTRICIONAL. Mantener seguimiento médico.",
			nil, nil,
			PriorityAttention,
			ColorBlue,
			MuacCodeFollow,
		),
	}
}
