package handler

import (
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	"github.com/go-playground/validator/v10"
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

func (h *SubscriptionPlanHandler) Create(c *fiber.Ctx) error {
	var req request.CreateSubscriptionPlanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrorResponse(err.Error(), nil))
	}
	userAuth := c.Locals("user").(middleware.AuthUser)
	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(errorMessages))
	}

	if err := h.subscriptionPlanUC.Create(req, userAuth); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Subscription Plan created successfully"))
}
