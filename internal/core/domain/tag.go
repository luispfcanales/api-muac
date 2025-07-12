// domain/tag.go
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Tag representa la entidad de etiqueta MUAC en el dominio
type Tag struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(100);not null;unique"`
	Description string    `json:"description" gorm:"column:description;type:text"`

	// Campos MUAC específicos
	Color    string `json:"color" gorm:"column:color;type:varchar(20)"`         // Código color hexadecimal
	Active   bool   `json:"active" gorm:"column:active;default:true"`           // Estado activo/inactivo
	MuacCode string `json:"muac_code" gorm:"column:muac_code;type:varchar(10)"` // MUAC-R1, MUAC-Y1, MUAC-G1
	Priority int    `json:"priority" gorm:"column:priority;type:int;default:1"` // 1-10 para ordenamiento

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Tag) TableName() string {
	return "tags"
}

// ============= CONSTRUCTORES =============

// NewTag crea una nueva instancia de Tag básica
func NewTag(name, description string) *Tag {
	return &Tag{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Active:      true,
		Priority:    PriorityLow,
		Color:       ColorGray,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewMuacTag crea una nueva etiqueta específica para MUAC
func NewMuacTag(name, description, color, muacCode string, priority int) *Tag {
	tag := &Tag{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Color:       color,
		Active:      true,
		MuacCode:    muacCode,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validar y usar color por defecto si es necesario
	if !IsValidHexColor(color) {
		tag.Color = tag.GetColorOrDefault()
	}

	return tag
}

// ============= VALIDACIÓN =============

// Validate valida que la etiqueta tenga los campos requeridos
func (t *Tag) Validate() error {
	if t.Name == "" {
		return ErrEmptyTagName
	}

	if t.MuacCode != "" && !IsValidMuacCode(t.MuacCode) {
		return fmt.Errorf("%w: %s", ErrInvalidMuacCode, t.MuacCode)
	}

	if t.Color != "" && !IsValidHexColor(t.Color) {
		return fmt.Errorf("%w: %s", ErrInvalidTagColor, t.Color)
	}

	if t.Priority < 1 || t.Priority > 10 {
		return ErrInvalidTagPriority
	}

	return nil
}

// ============= MÉTODOS DE ACTUALIZACIÓN =============

// Update actualiza los campos básicos de la etiqueta
func (t *Tag) Update(name, description string) {
	if name != "" {
		t.Name = name
	}
	if description != "" {
		t.Description = description
	}
	t.UpdatedAt = time.Now()
}

// UpdateMuacTag actualiza una etiqueta MUAC completa
func (t *Tag) UpdateMuacTag(name, description, color, muacCode string, priority int) error {
	// Validar antes de actualizar
	if name != "" {
		t.Name = name
	}
	if description != "" {
		t.Description = description
	}
	if color != "" {
		if !IsValidHexColor(color) {
			return fmt.Errorf("%w: %s", ErrInvalidTagColor, color)
		}
		t.Color = color
	}
	if muacCode != "" {
		if !IsValidMuacCode(muacCode) {
			return fmt.Errorf("%w: %s", ErrInvalidMuacCode, muacCode)
		}
		t.MuacCode = muacCode
	}
	if priority > 0 && priority <= 10 {
		t.Priority = priority
	}

	t.UpdatedAt = time.Now()
	return nil
}

// SetColor establece el color de la etiqueta
func (t *Tag) SetColor(color string) error {
	if !IsValidHexColor(color) {
		return fmt.Errorf("%w: %s", ErrInvalidTagColor, color)
	}
	t.Color = color
	t.UpdatedAt = time.Now()
	return nil
}

// SetMuacCode establece el código MUAC
func (t *Tag) SetMuacCode(muacCode string) error {
	if !IsValidMuacCode(muacCode) {
		return fmt.Errorf("%w: %s", ErrInvalidMuacCode, muacCode)
	}
	t.MuacCode = muacCode
	t.UpdatedAt = time.Now()
	return nil
}

// SetPriority establece la prioridad de la etiqueta
func (t *Tag) SetPriority(priority int) error {
	if priority < 1 || priority > 10 {
		return ErrInvalidTagPriority
	}
	t.Priority = priority
	t.UpdatedAt = time.Now()
	return nil
}

// Activate activa la etiqueta
func (t *Tag) Activate() {
	t.Active = true
	t.UpdatedAt = time.Now()
}

// Deactivate desactiva la etiqueta
func (t *Tag) Deactivate() {
	t.Active = false
	t.UpdatedAt = time.Now()
}

// ============= MÉTODOS DE CONSULTA =============

// IsActive verifica si la etiqueta está activa
func (t *Tag) IsActive() bool {
	return t.Active
}

// IsMuacTag verifica si es una etiqueta específica de MUAC
func (t *Tag) IsMuacTag() bool {
	return t.MuacCode != ""
}

// IsUrgent verifica si es una etiqueta de alta prioridad
func (t *Tag) IsUrgent() bool {
	return t.Priority >= PriorityExtreme || t.MuacCode == MuacCodeRed
}

// IsRisk verifica si es una etiqueta de riesgo
func (t *Tag) IsRisk() bool {
	return t.MuacCode == MuacCodeYellow
}

// IsNormal verifica si es una etiqueta de estado normal
func (t *Tag) IsNormal() bool {
	return t.MuacCode == MuacCodeGreen
}

// GetColorOrDefault retorna el color o un color por defecto
func (t *Tag) GetColorOrDefault() string {
	if t.Color != "" {
		return t.Color
	}

	switch t.MuacCode {
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

// GetPriorityText retorna el texto de prioridad
func (t *Tag) GetPriorityText() string {
	switch {
	case t.Priority >= PriorityCritical:
		return "Crítica"
	case t.Priority >= PriorityExtreme:
		return "Extrema"
	case t.Priority >= PriorityHigh:
		return "Alta"
	case t.Priority >= PriorityMedium:
		return "Media"
	default:
		return "Baja"
	}
}

// GetMuacDescription retorna la descripción basada en el código MUAC
func (t *Tag) GetMuacDescription() string {
	switch t.MuacCode {
	case MuacCodeRed:
		return fmt.Sprintf("Desnutrición aguda severa (SAM) - < %.1f cm - Requiere atención urgente", MuacThresholdSevere)
	case MuacCodeYellow:
		return fmt.Sprintf("Desnutrición aguda moderada (MAM) - %.1f-%.1f cm - Requiere seguimiento", MuacThresholdSevere, MuacThresholdModerate)
	case MuacCodeGreen:
		return fmt.Sprintf("Estado nutricional adecuado - ≥ %.1f cm - Mantener cuidados", MuacThresholdNormal)
	case MuacCodeFollow:
		return "Paciente en seguimiento post-intervención nutricional"
	default:
		return t.Description
	}
}

// ============= MÉTODOS ESTÁTICOS =============

// GetMuacTagForValue obtiene el tag MUAC apropiado para un valor
func GetMuacTagForValue(muacValue float64, tags []*Tag) *Tag {
	muacCode, colorCode, priority := ClassifyMuacValue(muacValue)

	// Buscar tag existente
	for _, tag := range tags {
		if tag.MuacCode == muacCode && tag.Active {
			return tag
		}
	}

	// Si no se encuentra, crear tag temporal
	return &Tag{
		Name:        getMuacTagName(muacCode),
		Description: getMuacTagDescription(muacCode),
		Color:       colorCode,
		MuacCode:    muacCode,
		Priority:    priority,
		Active:      true,
	}
}

// FilterActiveTags filtra solo tags activos
func FilterActiveTags(tags []*Tag) []*Tag {
	var activeTags []*Tag
	for _, tag := range tags {
		if tag.Active {
			activeTags = append(activeTags, tag)
		}
	}
	return activeTags
}

// SortTagsByPriority ordena tags por prioridad (mayor a menor)
func SortTagsByPriority(tags []*Tag) []*Tag {
	sorted := make([]*Tag, len(tags))
	copy(sorted, tags)

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

// ============= FUNCIONES HELPER PRIVADAS =============

// getMuacTagName retorna el nombre según código MUAC
func getMuacTagName(muacCode string) string {
	switch muacCode {
	case MuacCodeRed:
		return "ALERTA ROJA"
	case MuacCodeYellow:
		return "ALERTA AMARILLA"
	case MuacCodeGreen:
		return "ZONA VERDE"
	case MuacCodeFollow:
		return "SEGUIMIENTO"
	default:
		return "SIN CLASIFICAR"
	}
}

// getMuacTagDescription retorna la descripción según código MUAC
func getMuacTagDescription(muacCode string) string {
	switch muacCode {
	case MuacCodeRed:
		return fmt.Sprintf("Desnutrición aguda severa (SAM) - < %.1f cm - Requiere atención urgente", MuacThresholdSevere)
	case MuacCodeYellow:
		return fmt.Sprintf("Desnutrición aguda moderada (MAM) - %.1f-%.1f cm - Requiere seguimiento", MuacThresholdSevere, MuacThresholdModerate)
	case MuacCodeGreen:
		return fmt.Sprintf("Estado nutricional adecuado - ≥ %.1f cm - Mantener cuidados", MuacThresholdNormal)
	case MuacCodeFollow:
		return "Paciente en seguimiento post-intervención nutricional"
	default:
		return "Sin clasificación MUAC"
	}
}
