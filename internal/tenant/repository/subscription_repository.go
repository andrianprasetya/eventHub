package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(tx *gorm.DB, subscription *model.Subscription) error
}

type subscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{DB: db}
}

func (r subscriptionRepository) Create(tx *gorm.DB, subscription *model.Subscription) error {
	return tx.Create(&subscription).Error
}
