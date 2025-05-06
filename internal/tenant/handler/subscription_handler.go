package handler

import (
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	"github.com/gofiber/fiber/v2"
)

type SubscriptionPlanHandler struct {
	subscriptionPlanUC usecase.SubscriptionPlanUsecase
}

func NewSubscriptionPlanHandler(subscriptionPlanUC usecase.SubscriptionPlanUsecase) *SubscriptionPlanHandler {
	return &SubscriptionPlanHandler{subscriptionPlanUC: subscriptionPlanUC}
}

func (h *SubscriptionPlanHandler) GetAll(c *fiber.Ctx) error {
	subscriptionPlan, errRegisterTenant := h.subscriptionPlanUC.GetAll()
	if errRegisterTenant != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(errRegisterTenant.Error(), errRegisterTenant))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse("Get Subscription Plan successfully", subscriptionPlan))

}
