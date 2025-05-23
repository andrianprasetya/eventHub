package handler

import (
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	"github.com/andrianprasetya/eventHub/internal/ticket/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TicketHandler struct {
	ticketUC usecase.TicketUsecase
}

func NewTicketHandler(ticketUC usecase.TicketUsecase) *TicketHandler {
	return &TicketHandler{ticketUC: ticketUC}
}

func (u *TicketHandler) Create(c *fiber.Ctx) error {
	var req request.CreateTicketRequest
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

	if err := u.ticketUC.Create(req, userAuth, url); err != nil {
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
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(fiber.StatusCreated, "Create Ticket Successfully"))
}

func (u *TicketHandler) GetAll(c *fiber.Ctx) error {
	var query request.TicketPaginateParams

	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(http.StatusBadRequest, "invalid query parameters", err))
	}

	userAuth := c.Locals("user").(middleware.AuthUser)
	roles, total, err := u.ticketUC.GetAll(query, userAuth)
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
		"Get Ticket successfully",
		roles,
		query.Page,
		query.PageSize,
		total,
	))
}

func (u *TicketHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := u.ticketUC.GetByID(id)
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "get Ticket successfully", role))
}
