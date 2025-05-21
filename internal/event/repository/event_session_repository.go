package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventSessionRepository interface {
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, eventSession []*model.EventSession) error
}

type eventSessionRepository struct {
	DB *gorm.DB
}

func NewEventSessionRepository(db *gorm.DB) EventSessionRepository {
	return &eventSessionRepository{DB: db}
}

func (r *eventSessionRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, eventSession []*model.EventSession) error {
	return tx.WithContext(ctx).Create(eventSession).Error
}
