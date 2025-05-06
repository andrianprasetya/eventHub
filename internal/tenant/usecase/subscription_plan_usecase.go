package usecase

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
)

type SubscriptionPlanUsecase interface {
	GetAll() ([]*response.SubscriptionPlanListItemResponse, error)
}

type subscriptionPlanUsecase struct {
	subscriptionPlanRepo repository.SubscriptionPlanRepository
}

func NewSubscriptionPlanUsecase(subscriptionPlanRepo repository.SubscriptionPlanRepository) SubscriptionPlanUsecase {
	return &subscriptionPlanUsecase{subscriptionPlanRepo: subscriptionPlanRepo}
}

func (u subscriptionPlanUsecase) GetAll() ([]*response.SubscriptionPlanListItemResponse, error) {
	subscriptionPlan, err := u.subscriptionPlanRepo.GetAll()

	return mapper.FromSubscriptionPlanToList(subscriptionPlan), err
}
