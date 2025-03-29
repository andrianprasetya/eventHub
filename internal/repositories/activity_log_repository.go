package repositories

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	DB *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{DB: db}
}

func (r *ActivityLogRepository) Create(log *models.ActivityLog) error {
	return r.DB.Create(log).Error
}

func (r *ActivityLogRepository) GetByUser(userID string) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.DB.Where("user_id = ?", userID).Find(&logs).Error
	return logs, err
}
