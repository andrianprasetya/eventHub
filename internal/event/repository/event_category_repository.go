package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventCategoryRepository interface {
	Create(eventCategory *model.EventCategory) error
	CreateBulkWithTx(tx *gorm.DB, eventCategories *[]model.EventCategory) error
	AddCategoryToEventWithTx(tx *gorm.DB, id string, event *model.Event) error
	GetAll(page, pageSize int, tenantID *string) ([]*model.EventCategory, int64, error)
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

func (r *eventCategoryRepository) GetAll(page, pageSize int, tenantID *string) ([]*model.EventCategory, int64, error) {
	var eventCategories []*model.EventCategory
	var total int64
	db := r.DB.Model(&model.EventCategory{})

	if tenantID != nil {
		db = db.Where("tenant_id = ?", tenantID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&eventCategories).Error; err != nil {
		return nil, 0, err
	}
	return eventCategories, total, nil
}

func (r *eventCategoryRepository) AddCategoryToEventWithTx(tx *gorm.DB, id string, event *model.Event) error {
	return r.DB.Preload("Category").First(event, "id = ?", id).Error
}
