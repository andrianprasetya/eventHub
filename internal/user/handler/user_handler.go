package handler

import (
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUC usecase.UserUsecase
}

func NewUserHandler(userUC usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req request.CreateUserRequest
	var method = c.Method()
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrorResponse(err.Error(), nil))
	}

	userAuth := c.Locals("user").(middleware.AuthUser)

	if errValidation := validation.NewValidator().Validate(&req); errValidation != nil {
		errs := errValidation.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(errorMessages))
	}

	if errRegisterTenant := h.userUC.Create(req, userAuth, method); errRegisterTenant != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(errRegisterTenant.Error(), errRegisterTenant))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("User registered successfully"))
}
