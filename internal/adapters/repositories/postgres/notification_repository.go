package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// NotificationRepository implementa el repositorio de notificaciones usando PostgreSQL
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository crea una nueva instancia de NotificationRepository
func NewNotificationRepository(db *gorm.DB) ports.INotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

// Create crea una nueva notificación en la base de datos
func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetByID obtiene una notificación por su ID
func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	var notification domain.Notification
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&notification)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, result.Error
	}
	return &notification, nil
}

// GetAll obtiene todas las notificaciones
func (r *notificationRepository) GetAll(ctx context.Context) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	if err := r.db.WithContext(ctx).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// Update actualiza una notificación existente
func (r *notificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	result := r.db.WithContext(ctx).Save(notification)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

// Delete elimina una notificación por su ID
func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Notification{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}
