package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(tx *gorm.DB, event *model.Event) error
	GetAll(page, pageSize int) ([]*model.Event, int64, error)
	GetByID(id string) (*model.Event, error)
}

type eventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{DB: db}
}

func (r *eventRepository) Create(tx *gorm.DB, event *model.Event) error {
	return tx.Create(event).Error
}

func (r *eventRepository) GetAll(page, pageSize int) ([]*model.Event, int64, error) {
	var events []*model.Event
	var total int64
	if err := r.DB.Model(&model.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.DB.Limit(pageSize).Offset(offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func (r *eventRepository) GetByID(id string) (*model.Event, error) {
	var event model.Event
	if err := r.DB.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
