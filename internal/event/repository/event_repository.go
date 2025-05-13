package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(tx *gorm.DB, event *model.Event) error
}

type eventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{DB: db}
}

func (r eventRepository) Create(tx *gorm.DB, event *model.Event) error {
	return tx.Create(event).Error
}
