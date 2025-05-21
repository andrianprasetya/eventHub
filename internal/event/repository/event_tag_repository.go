package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"gorm.io/gorm"
)

type EventTagRepository interface {
	Create(ctx context.Context, eventTag *model.EventTag) error
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, eventTags []*model.EventTag) error
	GetAll(ctx context.Context, query request.EventTagPaginateRequest, tenantID *string) ([]*model.EventTag, int64, error)
}

type eventTagRepository struct {
	DB *gorm.DB
}

func NewEventTagRepository(db *gorm.DB) EventTagRepository {
	return &eventTagRepository{DB: db}
}

func (r *eventTagRepository) Create(ctx context.Context, eventTag *model.EventTag) error {
	return r.DB.WithContext(ctx).Create(eventTag).Error
}

func (r *eventTagRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, eventTags []*model.EventTag) error {
	return tx.WithContext(ctx).Create(eventTags).Error
}

func (r *eventTagRepository) GetAll(ctx context.Context, query request.EventTagPaginateRequest, tenantID *string) ([]*model.EventTag, int64, error) {
	var eventTags []*model.EventTag
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.EventTag{})

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

	if err := db.Limit(query.PageSize).Offset(offset).Find(&eventTags).Error; err != nil {
		return nil, 0, err
	}
	return eventTags, total, nil
}
