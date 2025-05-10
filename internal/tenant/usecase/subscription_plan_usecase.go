package usecase

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
	log "github.com/sirupsen/logrus"
)

type SubscriptionPlanUsecase interface {
	Create(req request.CreateSubscriptionPlanRequest, authUser middleware.AuthUser) error
	GetAll() ([]*response.SubscriptionPlanListItemResponse, error)
}

type subscriptionPlanUsecase struct {
	subscriptionPlanRepo repository.SubscriptionPlanRepository
}

func NewSubscriptionPlanUsecase(subscriptionPlanRepo repository.SubscriptionPlanRepository) SubscriptionPlanUsecase {
	return &subscriptionPlanUsecase{subscriptionPlanRepo: subscriptionPlanRepo}
}

func (u *subscriptionPlanUsecase) Create(req request.CreateSubscriptionPlanRequest, authUser middleware.AuthUser) error {
	features := utils.ToJSONString(req.Feature)

	subscriptionPlan := &model.SubscriptionPlan{
		ID:          utils.GenerateID(),
		Name:        req.Name,
		Price:       req.Price,
		Feature:     features,
		DurationDay: req.DurationDay,
	}

	if err := u.subscriptionPlanRepo.Create(subscriptionPlan); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create subscription plan")
		return fmt.Errorf("something Went wrong")
	}
	return nil
}

func (u *subscriptionPlanUsecase) GetAll() ([]*response.SubscriptionPlanListItemResponse, error) {
	subscriptionPlan, err := u.subscriptionPlanRepo.GetAll()

	return mapper.FromSubscriptionPlanToList(subscriptionPlan), err
}
