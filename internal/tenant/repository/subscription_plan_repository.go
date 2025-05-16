package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type SubscriptionPlanRepository interface {
	Create(subscriptionPlan *model.SubscriptionPlan) error
	GetAll(page, pageSize int) ([]*model.SubscriptionPlan, int64, error)
	GetByID(id string) (*model.SubscriptionPlan, error)
	Update(subscriptionPlan *model.SubscriptionPlan) error
	Delete(id string) error
}

type subscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{DB: db}
}

func (r *subscriptionPlanRepository) Create(subscriptionPlan *model.SubscriptionPlan) error {
	return r.DB.Create(subscriptionPlan).Error
}

func (r *subscriptionPlanRepository) GetAll(page, pageSize int) ([]*model.SubscriptionPlan, int64, error) {
	var subscriptionPlans []*model.SubscriptionPlan
	var total int64
	if err := r.DB.Model(&model.SubscriptionPlan{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.DB.Limit(pageSize).Offset(offset).Find(&subscriptionPlans).Error; err != nil {
		return nil, 0, err
	}
	return subscriptionPlans, total, nil
}

func (r *subscriptionPlanRepository) GetByID(id string) (*model.SubscriptionPlan, error) {
	var subscriptionPlan model.SubscriptionPlan
	if err := r.DB.First(&subscriptionPlan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subscriptionPlan, nil
}

func (r *subscriptionPlanRepository) Update(subscriptionPlan *model.SubscriptionPlan) error {
	return r.DB.Save(subscriptionPlan).Error
}

func (r *subscriptionPlanRepository) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&model.SubscriptionPlan{}).Error
}
