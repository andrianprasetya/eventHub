package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, discount []*model.Discount) error
}

type discountRepository struct {
	DB *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepository{DB: db}
}

func (r *discountRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, Discount []*model.Discount) error {
	return tx.WithContext(ctx).Create(Discount).Error
}
