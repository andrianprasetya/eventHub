package usecase

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
)

type SubscriptionPlanUsecase interface {
	Create(req request.CreateSubscriptionPlanRequest) error
	GetAll() ([]*response.SubscriptionPlanListItemResponse, error)
}

type subscriptionPlanUsecase struct {
	subscriptionPlanRepo repository.SubscriptionPlanRepository
}

func NewSubscriptionPlanUsecase(subscriptionPlanRepo repository.SubscriptionPlanRepository) SubscriptionPlanUsecase {
	return &subscriptionPlanUsecase{subscriptionPlanRepo: subscriptionPlanRepo}
}

func (u subscriptionPlanUsecase) Create(req request.CreateSubscriptionPlanRequest) error {
	
	return nil
}

func (u subscriptionPlanUsecase) GetAll() ([]*response.SubscriptionPlanListItemResponse, error) {
	subscriptionPlan, err := u.subscriptionPlanRepo.GetAll()

	return mapper.FromSubscriptionPlanToList(subscriptionPlan), err
}
