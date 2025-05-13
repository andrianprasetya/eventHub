package handler

import (
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TenantHandler struct {
	tenantUC usecase.TenantUsecase
}

func NewTenantHandler(tenantUC usecase.TenantUsecase) *TenantHandler {
	return &TenantHandler{tenantUC: tenantUC}
}

func (h *TenantHandler) RegisterTenant(c *fiber.Ctx) error {
	var req request.CreateTenantRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}

	if err := h.tenantUC.RegisterTenant(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(fiber.StatusCreated, "Tenant registered successfully"))
}

func (h *TenantHandler) UpdateTenant(c *fiber.Ctx) error {
	var req request.UpdateTenantRequest
	var id = c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}

	if err := h.tenantUC.Update(id, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, "Tenant updated successfully"))
}
