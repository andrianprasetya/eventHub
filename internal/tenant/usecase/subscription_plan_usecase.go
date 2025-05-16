package usecase

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
	log "github.com/sirupsen/logrus"
)

type SubscriptionPlanUsecase interface {
	Create(req request.CreateSubscriptionPlanRequest) (*response.SubscriptionPlanResponse, error)
	GetAll(page, pageSize int) ([]*response.SubscriptionPlanListItemResponse, int64, error)
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
		return &response.SubscriptionPlanResponse{}, fmt.Errorf("something Went wrong %w", err)
	}
	return mapper.FromSubscriptionModel(subscriptionPlan), nil
}

func (u *subscriptionPlanUsecase) GetAll(page, pageSize int) ([]*response.SubscriptionPlanListItemResponse, int64, error) {
	subscriptionPlan, total, err := u.subscriptionPlanRepo.GetAll(page, pageSize)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get subscription plan")
		return nil, 0, fmt.Errorf("something Went wrong %w", err)
	}

	return mapper.FromSubscriptionPlanToList(subscriptionPlan), total, err
}

func (u *subscriptionPlanUsecase) Update(id string, req request.UpdateSubscriptionPlanRequest) (*response.SubscriptionPlanResponse, error) {
	subscriptionPlan, err := u.subscriptionPlanRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get subscription plan")
		return &response.SubscriptionPlanResponse{}, fmt.Errorf("something Went wrong")
	}
	if req.Name != nil {
		subscriptionPlan.Name = *req.Name
	}
	if req.Price != nil {
		subscriptionPlan.Price = *req.Price
	}
	if req.Feature != nil {
		feature := utils.ToJSONString(*req.Feature)
		subscriptionPlan.Feature = feature
	}
	if req.DurationDay != nil {
		subscriptionPlan.DurationDay = *req.DurationDay
	}

	if err := u.subscriptionPlanRepo.Update(subscriptionPlan); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create subscription plan")
		return &response.SubscriptionPlanResponse{}, fmt.Errorf("something Went wrong")
	}

	return mapper.FromSubscriptionModel(subscriptionPlan), nil
}

func (u *subscriptionPlanUsecase) GetByID(id string) (*response.SubscriptionPlanResponse, error) {
	subscriptionPlan, err := u.subscriptionPlanRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get subscription plan")
		return &response.SubscriptionPlanResponse{}, fmt.Errorf("something Went wrong")
	}
	return mapper.FromSubscriptionModel(subscriptionPlan), nil
}

func (u *subscriptionPlanUsecase) Delete(id string) error {
	if err := u.subscriptionPlanRepo.Delete(id); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to delete subscription plan")
		return fmt.Errorf("something Went wrong")
	}
	return nil
}
