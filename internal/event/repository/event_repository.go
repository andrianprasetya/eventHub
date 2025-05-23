package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, tx *gorm.DB, event *model.Event) error
	GetAll(ctx context.Context, query request.EventPaginateRequest, tenantID *string) ([]*model.Event, int64, error)
	GetByID(ctx context.Context, id string) (*model.Event, error)
	CountCreatedEvent(ctx context.Context, tenantID string) (int, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, event *model.Event) error
	Update(ctx context.Context, event *model.Event) error
}

type eventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{DB: db}
}

func (r *eventRepository) Create(ctx context.Context, tx *gorm.DB, event *model.Event) error {
	return tx.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) GetAll(ctx context.Context, query request.EventPaginateRequest, tenantID *string) ([]*model.Event, int64, error) {
	var events []*model.Event
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.Event{})

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

func (r *eventRepository) GetByID(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event
	if err := r.DB.WithContext(ctx).First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, event *model.Event) error {
	return tx.WithContext(ctx).Save(event).Error
}

func (r *eventRepository) Update(ctx context.Context, event *model.Event) error {
	return r.DB.WithContext(ctx).Save(event).Error
}

func (r *eventRepository) CountCreatedEvent(ctx context.Context, tenantID string) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(model.Event{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return int(count), err
}
