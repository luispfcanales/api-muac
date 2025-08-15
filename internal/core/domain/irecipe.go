package domain

import (
	"time"

	"github.com/google/uuid"
)

// Recipe represents a nutritional recipe for children
type Recipe struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	MinAgeYears float32   `gorm:"type:float;not null" json:"min_age_years"`
	MaxAgeYears float32   `gorm:"type:float;not null" json:"max_age_years"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the name of the table in the database
func (Recipe) TableName() string {
	return "recipes"
}

// NewRecipe creates a new Recipe instance
func NewRecipe(title, content string, minAge, maxAge float32) *Recipe {
	return &Recipe{
		ID:          uuid.New(),
		Title:       title,
		Content:     content,
		MinAgeYears: minAge,
		MaxAgeYears: maxAge,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
