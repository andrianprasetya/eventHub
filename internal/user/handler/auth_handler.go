package handler

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userUC usecase.UserUsecase
}

func NewAuthHandler(userUC usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{userUC: userUC}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}

	ctx := context.Background()
	ip := c.IP()

	token, err := h.userUC.Login(ctx, req, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Get token successfully", token))

}
