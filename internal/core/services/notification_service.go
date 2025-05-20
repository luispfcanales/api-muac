package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// NotificationService implementa la lógica de negocio para notificaciones
type notificationService struct {
	notificationRepo ports.INotificationRepository
}

// NewNotificationService crea una nueva instancia de NotificationService
func NewNotificationService(notificationRepo ports.INotificationRepository) ports.INotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

// Create crea una nueva notificación
func (s *notificationService) Create(ctx context.Context, notification *domain.Notification) error {
	if err := notification.Validate(); err != nil {
		return err
	}
	return s.notificationRepo.Create(ctx, notification)
}

// GetByID obtiene una notificación por su ID
func (s *notificationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	return s.notificationRepo.GetByID(ctx, id)
}

// GetAll obtiene todas las notificaciones
func (s *notificationService) GetAll(ctx context.Context) ([]*domain.Notification, error) {
	return s.notificationRepo.GetAll(ctx)
}

// Update actualiza una notificación existente
func (s *notificationService) Update(ctx context.Context, notification *domain.Notification) error {
	if err := notification.Validate(); err != nil {
		return err
	}
	return s.notificationRepo.Update(ctx, notification)
}

// Delete elimina una notificación por su ID
func (s *notificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.Delete(ctx, id)
}
