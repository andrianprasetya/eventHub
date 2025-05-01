package handler

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/user"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userUC user.UserUsecase
}

func NewAuthHandler(userUC user.UserUsecase) *AuthHandler {
	return &AuthHandler{userUC: userUC}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrorResponse(err.Error(), nil))
	}

	if errValidation := validation.NewValidator().Validate(&req); errValidation != nil {
		errs := errValidation.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(errorMessages))
	}
	ctx := context.Background()
	token, errLogin := h.userUC.Login(ctx, req)
	if errLogin != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(errLogin.Error(), errLogin))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse("Get token successfully", token))

}
