package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventTagRepository interface {
	Create(eventTag *model.EventTag) error
	CreateBulkWithTx(tx *gorm.DB, eventTags *[]model.EventTag) error
	GetAll() ([]*model.EventTag, error)
}

type eventTagRepository struct {
	DB *gorm.DB
}

func NewEventTagRepository(db *gorm.DB) EventTagRepository {
	return &eventTagRepository{DB: db}
}

func (r *eventTagRepository) Create(eventTag *model.EventTag) error {
	return r.DB.Create(eventTag).Error
}

func (r *eventTagRepository) CreateBulkWithTx(tx *gorm.DB, eventTags *[]model.EventTag) error {
	return r.DB.Create(eventTags).Error
}

func (r *eventTagRepository) GetAll() ([]*model.EventTag, error) {
	var eventTags []*model.EventTag
	if err := r.DB.Find(&eventTags).Error; err != nil {
		return nil, err
	}
	return eventTags, nil
}
