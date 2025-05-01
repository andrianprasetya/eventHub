package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant"
	model "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) tenant.SubscriptionRepository {
	return &subscriptionRepository{DB: db}
}

func (r subscriptionRepository) Create(subscription *model.Subscription) error {
	return r.DB.Create(&subscription).Error
}
