package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
}

type eventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{DB: db}
}

func (r eventRepository) Create(event *model.Event) error {
	return r.DB.Create(event).Error
}
