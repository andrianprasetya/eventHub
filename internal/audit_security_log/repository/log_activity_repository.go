package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"gorm.io/gorm"
)

type LogActivityRepository interface {
	Create(ctx context.Context, history *model.ActivityLog) error
}

type logActivityRepository struct {
	DB *gorm.DB
}

func NewLogActivityRepository(db *gorm.DB) LogActivityRepository {
	return &logActivityRepository{DB: db}
}

func (r *logActivityRepository) Create(ctx context.Context, history *model.ActivityLog) error {
	return r.DB.WithContext(ctx).Create(history).Error
}
