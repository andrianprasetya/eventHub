package repositories

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.DB.Create(notification).Error
}

func (r *NotificationRepository) GetByID(id string) (*models.Notification, error) {
	var notification models.Notification
	err := r.DB.First(&notification, "id = ?", id).Error
	return &notification, err
}

func (r *NotificationRepository) GetUnreadByUser(userID string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.DB.Where("user_id = ? AND is_read = false", userID).Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(id string) error {
	return r.DB.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}
