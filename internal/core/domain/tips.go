package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tip represents a tip or advice entity
type Tip struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the name of the table in the database
func (Tip) TableName() string {
	return "tips"
}

// NewTip crea una nueva instancia de Tip básica
func NewTip(title, content string) *Tip {
	return &Tip{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
