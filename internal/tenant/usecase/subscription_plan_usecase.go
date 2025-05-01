package usecase

import (
	"github.com/andrianprasetya/eventHub/internal/tenant"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
)

type subscriptionPlanUsecase struct {
	subscriptionPlanRepo tenant.SubscriptionPlanRepository
}

func NewSubscriptionPlanUsecase(subscriptionPlanRepo tenant.SubscriptionPlanRepository) tenant.SubscriptionPlanUsecase {
	return &subscriptionPlanUsecase{subscriptionPlanRepo: subscriptionPlanRepo}
}

func (u subscriptionPlanUsecase) GetAll() ([]*response.SubscriptionPlanListItemResponse, error) {
	subscriptionPlan, err := u.subscriptionPlanRepo.GetAll()

	return mapper.FromSubscriptionPlanToList(subscriptionPlan), err
}
