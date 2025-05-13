package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventSessionRepository interface {
	CreateBulkWithTx(tx *gorm.DB, eventSession []*model.EventSession) error
}

type eventSessionRepository struct {
	DB *gorm.DB
}

func NewEventSessionRepository(db *gorm.DB) EventSessionRepository {
	return &eventSessionRepository{DB: db}
}

func (r *eventSessionRepository) CreateBulkWithTx(tx *gorm.DB, eventSession []*model.EventSession) error {
	return tx.Create(eventSession).Error
}
