package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, discounts []*model.Discount) error
	CreateBulk(ctx context.Context, discounts []*model.Discount) error
	GetAll(ctx context.Context, query request.DiscountPaginateParams, tenantID *string) ([]*model.Discount, int64, error)
	GetByID(ctx context.Context, id string) (*model.Discount, error)
}

type discountRepository struct {
	DB *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepository{DB: db}
}

func (r *discountRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, discounts []*model.Discount) error {
	return tx.WithContext(ctx).Create(discounts).Error
}

func (r *discountRepository) CreateBulk(ctx context.Context, discounts []*model.Discount) error {
	return r.DB.WithContext(ctx).Create(discounts).Error
}

func (r *discountRepository) GetAll(ctx context.Context, query request.DiscountPaginateParams, tenantID *string) ([]*model.Discount, int64, error) {
	var discounts []*model.Discount
	var total int64

	db := r.DB.WithContext(ctx).Preload("Event").Model(&model.Discount{})
	if tenantID != nil {
		db = db.Joins("JOIN events ON events.id = discounts.event_id").Where("events.tenant_id = ?", *tenantID)
	}

	if query.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*query.Name+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize

	if err := db.Limit(query.PageSize).Offset(offset).Find(&discounts).Error; err != nil {
		return nil, 0, err
	}
	return discounts, total, nil
}

func (r *discountRepository) GetByID(ctx context.Context, id string) (*model.Discount, error) {
	var discount model.Discount
	if err := r.DB.WithContext(ctx).First(&discount, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &discount, nil
}
