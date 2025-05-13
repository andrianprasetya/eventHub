package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventCategoryRepository interface {
	Create(eventCategory *model.EventCategory) error
	CreateBulkWithTx(tx *gorm.DB, eventCategories *[]model.EventCategory) error
	AddCategoryToEventWithTx(tx *gorm.DB, id string, event *model.Event) error
	GetAll() ([]*model.EventCategory, error)
}

type eventCategoryRepository struct {
	DB *gorm.DB
}

func NewEventCategoryRepository(db *gorm.DB) EventCategoryRepository {
	return &eventCategoryRepository{DB: db}
}

func (r *eventCategoryRepository) Create(eventCategory *model.EventCategory) error {
	return r.DB.Create(eventCategory).Error
}

func (r *eventCategoryRepository) CreateBulkWithTx(tx *gorm.DB, eventCategories *[]model.EventCategory) error {
	return r.DB.Create(eventCategories).Error
}

func (r *eventCategoryRepository) GetAll() ([]*model.EventCategory, error) {
	var eventCategories []*model.EventCategory
	if err := r.DB.Find(&eventCategories).Error; err != nil {
		return nil, err
	}
	return eventCategories, nil
}

func (r *eventCategoryRepository) AddCategoryToEventWithTx(tx *gorm.DB, id string, event *model.Event) error {
	return tx.Preload("Category").First(event, "id = ?", id).Error
}
