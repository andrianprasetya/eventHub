package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventCategoryRepository interface {
	Create(eventCategory *model.EventCategory) error
	CreateBulkWithTx(tx *gorm.DB, eventCategories *[]model.EventCategory) error
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
