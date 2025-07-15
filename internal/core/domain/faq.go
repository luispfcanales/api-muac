package domain

import (
	"time"

	"github.com/google/uuid"
)

// Constantes para categorías de FAQs
const (
	FAQCategoryTapeAndApp    = "SOBRE EL USO DE LA CINTA Y EL APP"
	FAQCategoryAppInfo       = "SOBRE EL FUNCIONAMIENTO DEL APLICATIVO"
	FAQCategoryResults       = "SOBRE LOS RESULTADOS Y LO QUE DEBO HACER"
	FAQCategoryHealthCenters = "SOBRE LOS CENTROS DE SALUD Y EL APOYO LOCAL"
	FAQCategoryPrivacy       = "SOBRE PRIVACIDAD Y SEGURIDAD"
	FAQCategoryOther         = "OTRAS PREGUNTAS"
)

// Lista de todas las categorías válidas
var ValidFAQCategories = []string{
	FAQCategoryTapeAndApp,
	FAQCategoryResults,
	FAQCategoryHealthCenters,
	FAQCategoryPrivacy,
	FAQCategoryOther,
}

// FAQGrouped representa FAQs agrupadas por categoría
type FAQGrouped struct {
	Category string `json:"category"`
	FAQs     []*FAQ `json:"faqs"`
}

// FAQ representa la entidad de pregunta frecuente en el dominio
type FAQ struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Question  string    `json:"question" gorm:"column:question;type:text;not null"`
	Answer    string    `json:"answer" gorm:"column:answer;type:text;not null"`
	Category  string    `json:"category" gorm:"column:category;type:varchar(100);not null;default:'OTRAS PREGUNTAS'"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (FAQ) TableName() string {
	return "faqs"
}

// NewFAQ crea una nueva instancia de FAQ
func NewFAQ(question, answer, category string) (*FAQ, error) {
	if category == "" {
		category = FAQCategoryOther
	}
	return &FAQ{
		ID:        uuid.New(),
		Question:  question,
		Answer:    answer,
		Category:  category,
		CreatedAt: time.Now(),
	}, nil
}

// Validate valida que la FAQ tenga los campos requeridos
func (f *FAQ) Validate() error {
	if f.Question == "" {
		return ErrEmptyFAQQuestion
	}
	if f.Answer == "" {
		return ErrEmptyFAQAnswer
	}

	// Validar que la categoría sea una de las permitidas
	valid := false
	for _, cat := range ValidFAQCategories {
		if cat == f.Category {
			valid = true
			break
		}
	}

	if !valid {
		return ErrInvalidFAQCategory
	}

	return nil
}

// Update actualiza los campos de la FAQ
func (f *FAQ) Update(question, answer, category string) error {
	if question != "" {
		f.Question = question
	}
	if answer != "" {
		f.Answer = answer
	}
	if category != "" {
		f.Category = category
	}
	f.UpdatedAt = time.Now()
	return nil
}
