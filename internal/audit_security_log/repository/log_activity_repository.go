package repository

import (
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"gorm.io/gorm"
)

type LogActivityRepository interface {
	Create(history *model.ActivityLog) error
}

type logActivityRepository struct {
	DB *gorm.DB
}

func NewLogActivityRepository(db *gorm.DB) LogActivityRepository {
	return &logActivityRepository{DB: db}
}

func (r *logActivityRepository) Create(history *model.ActivityLog) error {
	return r.DB.Create(history).Error
}
