package domain

import (
	"time"

	"github.com/google/uuid"
)

// Measurement representa la entidad de medición en el dominio
type Measurement struct {
	ID               uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey"`
	MuacValue        float64         `json:"muac_value" gorm:"column:muac_value;type:decimal(10,2);not null"`
	Description      string          `json:"description" gorm:"column:description;type:text"`
	PatientID        uuid.UUID       `json:"patient_id" gorm:"column:patient_id;type:uuid;not null"`
	UserID           uuid.UUID       `json:"user_id" gorm:"column:user_id;type:uuid;not null"`
	TagID            *uuid.UUID      `json:"tag_id,omitempty" gorm:"column:tag_id;type:uuid"`
	RecommendationID *uuid.UUID      `json:"recommendation_id,omitempty" gorm:"column:recommendation_id;type:uuid"`
	CreatedAt        time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Patient          *Patient        `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	User             *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tag              *Tag            `json:"tag,omitempty" gorm:"foreignKey:TagID"`
	Recommendation   *Recommendation `json:"recommendation" gorm:"foreignKey:RecommendationID"`
}

// TableName especifica el nombre de la tabla para GORM
func (Measurement) TableName() string {
	return "measurements"
}

// NewMeasurement crea una nueva instancia de Measurement
func NewMeasurement(muacValue float64, description string, timestamp time.Time, patientID, userID uuid.UUID, tagID, recommendationID *uuid.UUID) *Measurement {
	if tagID != nil && *tagID == uuid.Nil {
		tagID = nil
	}
	if recommendationID != nil && *recommendationID == uuid.Nil {
		recommendationID = nil
	}
	return &Measurement{
		ID:          uuid.New(),
		MuacValue:   muacValue,
		Description: description,
		PatientID:   patientID,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}
}

// Validate valida que la medición tenga los campos requeridos
func (m *Measurement) Validate() error {
	if m.MuacValue <= 0 {
		return ErrInvalidMuacValue
	}
	if m.PatientID == uuid.Nil {
		return ErrEmptyPatientID
	}
	if m.UserID == uuid.Nil {
		return ErrEmptyUserID
	}
	return nil
}

// Update actualiza los campos de la medición
func (m *Measurement) Update(muacValue float64, description, location string, timestamp time.Time, tagID, recommendationID *uuid.UUID) {
	m.MuacValue = muacValue
	m.Description = description
	m.TagID = tagID
	m.RecommendationID = recommendationID
	m.UpdatedAt = time.Now()
}

// SetTag asigna una etiqueta a la medición
func (m *Measurement) SetTag(tagID *uuid.UUID) {
	m.TagID = tagID
	m.UpdatedAt = time.Now()
}

// SetRecommendation asigna una recomendación a la medición
func (m *Measurement) SetRecommendation(recommendationID *uuid.UUID) {
	m.RecommendationID = recommendationID
	m.UpdatedAt = time.Now()
}
