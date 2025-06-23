package domain

import (
	"time"

	"github.com/google/uuid"
)

// Notification representa la entidad de notificación en el dominio
type Notification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title     string    `json:"title" gorm:"column:TITLE;type:varchar(255);not null"`
	Body      string    `json:"body" gorm:"column:BODY;type:text"`
	Visible   bool      `json:"visible" gorm:"column:VISIBLE;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Notification) TableName() string {
	return "notifications"
}

// NewNotification crea una nueva instancia de Notification
func NewNotification(title, body string, visible bool) *Notification {
	return &Notification{
		ID:        uuid.New(),
		Title:     title,
		Body:      body,
		Visible:   visible,
		CreatedAt: time.Now(),
	}
}

// Validate valida que la notificación tenga los campos requeridos
func (n *Notification) Validate() error {
	if n.Title == "" {
		return ErrEmptyNotificationTitle
	}
	return nil
}

// Update actualiza los campos de la notificación
func (n *Notification) Update(title, body string, visible bool) {
	n.Title = title
	n.Body = body
	n.Visible = visible
	n.UpdatedAt = time.Now()
}

// SetVisible establece la visibilidad de la notificación
func (n *Notification) SetVisible(visible bool) {
	n.Visible = visible
	n.UpdatedAt = time.Now()
}
