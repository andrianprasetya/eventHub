package handler

import (
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type SubscriptionPlanHandler struct {
	subscriptionPlanUC usecase.SubscriptionPlanUsecase
}

func NewSubscriptionPlanHandler(subscriptionPlanUC usecase.SubscriptionPlanUsecase) *SubscriptionPlanHandler {
	return &SubscriptionPlanHandler{subscriptionPlanUC: subscriptionPlanUC}
}

func (h *SubscriptionPlanHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	subscriptionPlan, total, err := h.subscriptionPlanUC.GetAll(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(
		fiber.StatusOK,
		"Get Subscription Plan successfully",
		subscriptionPlan,
		page,
		pageSize,
		total,
	))
}

func (h *SubscriptionPlanHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	subscriptionPlan, err := h.subscriptionPlanUC.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Get Subscription Plan successfully", subscriptionPlan))
}

func (h *SubscriptionPlanHandler) Create(c *fiber.Ctx) error {
	var req request.CreateSubscriptionPlanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}
	userAuth := c.Locals("user").(middleware.AuthUser)
	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}
	subscriptionPlan, err := h.subscriptionPlanUC.Create(req, userAuth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Create Subscription Plan successfully", subscriptionPlan))
}

func (h *SubscriptionPlanHandler) Update(c *fiber.Ctx) error {
	var req request.UpdateSubscriptionPlanRequest
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}

	subscriptionPlan, err := h.subscriptionPlanUC.Update(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Create Subscription Plan successfully", subscriptionPlan))
}

func (h *SubscriptionPlanHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.subscriptionPlanUC.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, "Create Subscription Plan successfully"))
}
