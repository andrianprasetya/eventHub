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
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrorResponse(err.Error(), nil))
	}

	if errValidation := validation.NewValidator().Validate(&req); errValidation != nil {
		errs := errValidation.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(errorMessages))
	}

	if errRegisterTenant := h.tenantUC.RegisterTenant(req); errRegisterTenant != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(errRegisterTenant.Error(), errRegisterTenant))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Tenant registered successfully"))

}
