package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventCategoryRepository interface {
	Create(ctx context.Context, eventCategory *model.EventCategory) error
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, eventCategories []*model.EventCategory) error
	AddCategoryToEventWithTx(ctx context.Context, tx *gorm.DB, event *model.Event) error
	GetAll(ctx context.Context, query request.EventCategoryPaginateRequest, tenantID *string) ([]*model.EventCategory, int64, error)
}

type eventCategoryRepository struct {
	DB *gorm.DB
}

func NewEventCategoryRepository(db *gorm.DB) EventCategoryRepository {
	return &eventCategoryRepository{DB: db}
}

func (r *eventCategoryRepository) Create(ctx context.Context, eventCategory *model.EventCategory) error {
	return r.DB.WithContext(ctx).Create(eventCategory).Error
}

func (r *eventCategoryRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, eventCategories []*model.EventCategory) error {
	return tx.WithContext(ctx).Create(eventCategories).Error
}

func (r *eventCategoryRepository) GetAll(ctx context.Context, query request.EventCategoryPaginateRequest, tenantID *string) ([]*model.EventCategory, int64, error) {
	var eventCategories []*model.EventCategory
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.EventCategory{})

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
	if err := db.Limit(query.PageSize).Offset(offset).Find(&eventCategories).Error; err != nil {
		return nil, 0, err
	}
	return eventCategories, total, nil
}

func (r *eventCategoryRepository) AddCategoryToEventWithTx(ctx context.Context, tx *gorm.DB, event *model.Event) error {
	return tx.WithContext(ctx).Preload("Category").First(event, "id = ?", event.ID).Error
}
