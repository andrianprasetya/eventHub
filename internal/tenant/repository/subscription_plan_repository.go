package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type SubscriptionPlanRepository interface {
	Create(ctx context.Context, subscriptionPlan *model.SubscriptionPlan) error
	GetAll(ctx context.Context, query request.SubscriptionPaginateParams) ([]*model.SubscriptionPlan, int64, error)
	GetByID(ctx context.Context, id string) (*model.SubscriptionPlan, error)
	Update(subscriptionPlan *model.SubscriptionPlan) error
	Delete(ctx context.Context, id string) error
}

type subscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{DB: db}
}

func (r *subscriptionPlanRepository) Create(ctx context.Context, subscriptionPlan *model.SubscriptionPlan) error {
	return r.DB.WithContext(ctx).Create(subscriptionPlan).Error
}

func (r *subscriptionPlanRepository) GetAll(ctx context.Context, query request.SubscriptionPaginateParams) ([]*model.SubscriptionPlan, int64, error) {
	var subscriptionPlans []*model.SubscriptionPlan
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.SubscriptionPlan{})

	if query.Name != nil {
		db = db.Where("name = ?", "%"+*query.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize

	if err := db.Limit(query.PageSize).Offset(offset).Find(&subscriptionPlans).Error; err != nil {
		return nil, 0, err
	}
	return subscriptionPlans, total, nil
}

func (r *subscriptionPlanRepository) GetByID(ctx context.Context, id string) (*model.SubscriptionPlan, error) {
	var subscriptionPlan model.SubscriptionPlan
	if err := r.DB.WithContext(ctx).First(&subscriptionPlan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subscriptionPlan, nil
}

func (r *subscriptionPlanRepository) Update(subscriptionPlan *model.SubscriptionPlan) error {
	return r.DB.Save(subscriptionPlan).Error
}

func (r *subscriptionPlanRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.SubscriptionPlan{}).Error
}
