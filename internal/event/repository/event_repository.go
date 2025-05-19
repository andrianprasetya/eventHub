package repository

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(tx *gorm.DB, event *model.Event) error
	GetAll(query request.EventPaginateRequest, tenantID *string) ([]*model.Event, int64, error)
	GetByID(id string) (*model.Event, error)
	CountCreatedEvent(tenantID string) int
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

func (r *eventRepository) GetAll(query request.EventPaginateRequest, tenantID *string) ([]*model.Event, int64, error) {
	var events []*model.Event
	var total int64

	db := r.DB.Model(&model.Event{})

	if query.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*query.Name+"%")
	}
	if tenantID != nil {
		db = db.Where("tenant_id = ?", tenantID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize

	if err := r.DB.Limit(query.PageSize).Offset(offset).Find(&events).Error; err != nil {
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

func (r *eventRepository) CountCreatedEvent(tenantID string) int {
	var count int64
	if err := r.DB.Model(model.Event{}).Where("tenant_id = ?", tenantID).Count(&count).Error; err != nil {
		return 0
	}
	return int(count)
}
