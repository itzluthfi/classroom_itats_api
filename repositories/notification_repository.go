package repositories

import (
	"classroom_itats_api/entities"
	"context"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Save(ctx context.Context, n *entities.Notification) error
	GetByRecipient(ctx context.Context, recipientID string, limit int) ([]entities.Notification, error)
	MarkOneRead(ctx context.Context, id int64, recipientID string) error
	MarkAllRead(ctx context.Context, recipientID string) error
	UnreadCount(ctx context.Context, recipientID string) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Save(ctx context.Context, n *entities.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationRepository) GetByRecipient(ctx context.Context, recipientID string, limit int) ([]entities.Notification, error) {
	var notifs []entities.Notification
	if limit <= 0 {
		limit = 50
	}
	err := r.db.WithContext(ctx).
		Where("recipient_id = ?", recipientID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) MarkOneRead(ctx context.Context, id int64, recipientID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("id = ? AND recipient_id = ?", id, recipientID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, recipientID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("recipient_id = ? AND is_read = false", recipientID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) UnreadCount(ctx context.Context, recipientID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("recipient_id = ? AND is_read = false", recipientID).
		Count(&count).Error
	return count, err
}
