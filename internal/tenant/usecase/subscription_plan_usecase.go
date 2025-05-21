package usecase

import (
	"context"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

type SubscriptionPlanUsecase interface {
	Create(req request.CreateSubscriptionPlanRequest) (*response.SubscriptionPlanResponse, error)
	GetAll(query request.SubscriptionPaginateParams) ([]*response.SubscriptionPlanListItemResponse, int64, error)
	GetByID(id string) (*response.SubscriptionPlanResponse, error)
	Update(id string, req request.UpdateSubscriptionPlanRequest) (*response.SubscriptionPlanResponse, error)
	Delete(id string) error
}

type subscriptionPlanUsecase struct {
	subscriptionPlanRepo repository.SubscriptionPlanRepository
}

func NewSubscriptionPlanUsecase(subscriptionPlanRepo repository.SubscriptionPlanRepository) SubscriptionPlanUsecase {
	return &subscriptionPlanUsecase{subscriptionPlanRepo: subscriptionPlanRepo}
}

func (u *subscriptionPlanUsecase) Create(req request.CreateSubscriptionPlanRequest) (*response.SubscriptionPlanResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	features, err := utils.ToJSONString(req.Feature)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to modify Json to string")
		return nil, appErrors.ErrInternalServer
	}
	subscriptionPlan := service.MapSubscriptionPlanPayload(req, features)

	if err := u.subscriptionPlanRepo.Create(ctx, subscriptionPlan); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create subscription plan")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromSubscriptionModel(subscriptionPlan), nil
}

func (u *subscriptionPlanUsecase) GetAll(query request.SubscriptionPaginateParams) ([]*response.SubscriptionPlanListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	subscriptionPlan, total, err := u.subscriptionPlanRepo.GetAll(ctx, query)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get subscription plan")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromSubscriptionPlanToList(subscriptionPlan), total, err
}

func (u *subscriptionPlanUsecase) Update(id string, req request.UpdateSubscriptionPlanRequest) (*response.SubscriptionPlanResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	subscriptionPlan, err := u.subscriptionPlanRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get subscription plan")
		return nil, appErrors.ErrInternalServer
	}
	if req.Name != nil {
		subscriptionPlan.Name = *req.Name
	}
	if req.Price != nil {
		subscriptionPlan.Price = *req.Price
	}
	if req.Feature != nil {
		featureUpdated, err := utils.ToJSONString(*req.Feature)
		if err != nil {
			log.WithFields(log.Fields{
				"errors": err,
			}).Error("failed to modify Json to string")
			return nil, appErrors.ErrInternalServer
		}
		subscriptionPlan.Feature = featureUpdated
	}
	if req.DurationDay != nil {
		subscriptionPlan.DurationDay = *req.DurationDay
	}

	if err := u.subscriptionPlanRepo.Update(subscriptionPlan); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create subscription plan")
		return nil, appErrors.ErrInternalServer
	}

	return mapper.FromSubscriptionModel(subscriptionPlan), err
}

func (u *subscriptionPlanUsecase) GetByID(id string) (*response.SubscriptionPlanResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	subscriptionPlan, err := u.subscriptionPlanRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get subscription plan")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromSubscriptionModel(subscriptionPlan), err
}

func (u *subscriptionPlanUsecase) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := u.subscriptionPlanRepo.Delete(ctx, id); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to delete subscription plan")
		return appErrors.ErrInternalServer
	}
	return nil
}
