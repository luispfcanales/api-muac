package domain

import (
	"time"

	"github.com/google/uuid"
)

// FAQ representa la entidad de pregunta frecuente en el dominio
type FAQ struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Question  string    `json:"question" gorm:"column:QUESTION;type:text;not null"`
	Answer    string    `json:"answer" gorm:"column:ANSWER;type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (FAQ) TableName() string {
	return "FAQ"
}

// NewFAQ crea una nueva instancia de FAQ
func NewFAQ(question, answer string) *FAQ {
	return &FAQ{
		ID:        uuid.New(),
		Question:  question,
		Answer:    answer,
		CreatedAt: time.Now(),
	}
}

// Validate valida que la FAQ tenga los campos requeridos
func (f *FAQ) Validate() error {
	if f.Question == "" {
		return ErrEmptyFAQQuestion
	}
	if f.Answer == "" {
		return ErrEmptyFAQAnswer
	}
	return nil
}

// Update actualiza los campos de la FAQ
func (f *FAQ) Update(question, answer string) {
	f.Question = question
	f.Answer = answer
	f.UpdatedAt = time.Now()
}