package service

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
)

func MapSubscriptionPlanPayload(req request.CreateSubscriptionPlanRequest, features string) *model.SubscriptionPlan {
	return &model.SubscriptionPlan{
		ID:          utils.GenerateID(),
		Name:        req.Name,
		Price:       req.Price,
		Feature:     features,
		DurationDay: req.DurationDay,
	}
}
