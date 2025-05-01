package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type subscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) tenant.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{DB: db}
}

func (r subscriptionPlanRepository) GetAll() ([]*model.SubscriptionPlan, error) {
	var subscriptionPlans []*model.SubscriptionPlan
	if err := r.DB.Find(&subscriptionPlans).Error; err != nil {
		return nil, err
	}
	return subscriptionPlans, nil
}

func (r subscriptionPlanRepository) Get(id string) (*model.SubscriptionPlan, error) {
	var subscriptionPlan model.SubscriptionPlan
	if err := r.DB.First(&subscriptionPlan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subscriptionPlan, nil
}
