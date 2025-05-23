package handler

import (
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type UserHandler struct {
	userUC usecase.UserUsecase
}

func NewUserHandler(userUC usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req request.CreateUserRequest
	var url = c.OriginalURL()
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	userAuth := c.Locals("user").(middleware.AuthUser)

	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}

	if err := h.userUC.Create(req, &userAuth, url); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			message := appErr.Message
			var errRes error
			if appErr.ShouldExpose() {
				errRes = appErr.Err
			}
			return c.Status(appErr.StatusCode()).JSON(response.ErrorResponse(appErr.StatusCode(), message, errRes))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, "User registered successfully"))
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	var query request.UserPaginateParams

	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(http.StatusBadRequest, "invalid query parameters", err))
	}

	userAuth := c.Locals("user").(middleware.AuthUser)

	var tenantID *string
	if userAuth.Tenant.ID != "" {
		tenantID = &userAuth.Tenant.ID
	}

	users, total, err := h.userUC.GetAll(query, tenantID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			message := appErr.Message
			var errRes error
			if appErr.ShouldExpose() {
				errRes = appErr.Err
			}
			return c.Status(appErr.StatusCode()).JSON(response.ErrorResponse(appErr.StatusCode(), message, errRes))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(
		fiber.StatusOK,
		"Get Users Successfully",
		users, query.Page,
		query.PageSize, total))
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userUC.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			message := appErr.Message
			var errRes error
			if appErr.ShouldExpose() {
				errRes = appErr.Err
			}
			return c.Status(appErr.StatusCode()).JSON(response.ErrorResponse(appErr.StatusCode(), message, errRes))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Get User Successfully", user))
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var req request.UpdateUserRequest
	id := c.Params("id")
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	if err := validation.NewValidator().Validate(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(errs)
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}

	if err := h.userUC.Update(req, id); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			message := appErr.Message
			var errRes error
			if appErr.ShouldExpose() {
				errRes = appErr.Err
			}
			return c.Status(appErr.StatusCode()).JSON(response.ErrorResponse(appErr.StatusCode(), message, errRes))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, "Update User Successfully"))

}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userUC.Delete(id); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			message := appErr.Message
			var errRes error
			if appErr.ShouldExpose() {
				errRes = appErr.Err
			}
			return c.Status(appErr.StatusCode()).JSON(response.ErrorResponse(appErr.StatusCode(), message, errRes))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, "Delete User Successfully"))
}
