package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, subscription *model.Subscription) error
}

type subscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{DB: db}
}

func (r subscriptionRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, subscription *model.Subscription) error {
	return tx.WithContext(ctx).Create(&subscription).Error
}
