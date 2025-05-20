package domain

import (
	"time"

	"github.com/google/uuid"
)

// Recommendation representa la entidad de recomendación médica en el dominio
type Recommendation struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name               string    `json:"name" gorm:"column:NAME;type:varchar(100);not null"`
	Description        string    `json:"description" gorm:"column:DESCRIPTION;type:text"`
	RecommendationUmbral string  `json:"recommendation_umbral" gorm:"column:RECOMMENDATION_UMBRAL;type:varchar(255)"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Recommendation) TableName() string {
	return "RECOMMENDATION"
}

// NewRecommendation crea una nueva instancia de Recommendation
func NewRecommendation(name, description, umbral string) *Recommendation {
	return &Recommendation{
		ID:                 uuid.New(),
		Name:               name,
		Description:        description,
		RecommendationUmbral: umbral,
		CreatedAt:          time.Now(),
	}
}

// Validate valida que la recomendación tenga los campos requeridos
func (r *Recommendation) Validate() error {
	if r.Name == "" {
		return ErrEmptyRecommendationName
	}
	return nil
}

// Update actualiza los campos de la recomendación
func (r *Recommendation) Update(name, description, umbral string) {
	r.Name = name
	r.Description = description
	r.RecommendationUmbral = umbral
	r.UpdatedAt = time.Now()
}