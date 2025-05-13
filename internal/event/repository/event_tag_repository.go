package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventTagRepository interface {
	Create(eventTag *model.EventTag) error
	CreateBulkWithTx(tx *gorm.DB, eventTags *[]model.EventTag) error
	GetAll(page, pageSize int) ([]*model.EventTag, int64, error)
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

func (r *eventTagRepository) GetAll(page, pageSize int) ([]*model.EventTag, int64, error) {
	var eventTags []*model.EventTag
	var total int64

	if err := r.DB.Model(&model.EventTag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize

	if err := r.DB.Limit(pageSize).Offset(offset).Find(&eventTags).Error; err != nil {
		return nil, 0, err
	}
	return eventTags, total, nil
}
