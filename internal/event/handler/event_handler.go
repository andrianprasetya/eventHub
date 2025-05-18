package handler

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/usecase"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type EventHandler struct {
	eventUC usecase.EventUsecase
}

func NewEventHandler(eventUC usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUC: eventUC}
}

func (h *EventHandler) GetTags(c *fiber.Ctx) error {
	var query request.EventTagPaginateRequest
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(http.StatusBadRequest, "invalid query parameters", err))
	}
	userAuth := c.Locals("user").(middleware.AuthUser)

	eventTags, total, err := h.eventUC.GetTags(query, &userAuth.Tenant.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(fiber.StatusOK, "Get Event Tags successfully", eventTags, query.Page, query.PageSize, total))
}

func (h *EventHandler) GetCategories(c *fiber.Ctx) error {
	var query request.EventCategoryPaginateRequest
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(http.StatusBadRequest, "invalid query parameters", err))
	}
	userAuth := c.Locals("user").(middleware.AuthUser)

	eventCategories, total, err := h.eventUC.GetCategories(query, &userAuth.Tenant.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(fiber.StatusOK, "Get Event Categories successfully", eventCategories, query.Page, query.PageSize, total))
}

func (h *EventHandler) Create(c *fiber.Ctx) error {
	var req request.CreateEventRequest

	var url = c.OriginalURL()

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationResponse(fiber.StatusUnprocessableEntity, err))
	}

	userAuth := c.Locals("user").(middleware.AuthUser)

	if errValidation := validation.NewValidator().Validate(&req); errValidation != nil {
		errs := errValidation.(validator.ValidationErrors)
		errorMessages := validation.MapValidationErrorsToJSONTags(req, errs)

		if req.EndDate.Before(req.StartDate) {
			errorMessages["end_date"] = "end_date must not below start_date"
		}
		for i, discount := range req.Discounts {
			if discount.EndDate.Before(discount.StartDate) {
				key := fmt.Sprintf("discounts[%d].end_date", i)
				errorMessages[key] = "end_date must not below start_date"
			}
			if discount.StartDate.Before(req.StartDate) {
				key := fmt.Sprintf("discounts[%d].start_date", i)
				errorMessages[key] = "start_date discount must not below start_date event"
			}
			if discount.EndDate.After(req.EndDate) {
				key := fmt.Sprintf("discounts[%d].end_date", i)
				errorMessages[key] = "end_date discount must not above end_start event"
			}
		}
		for i, session := range req.Sessions {
			if session.EndDateTime.Before(session.StartDateTime) {
				key := fmt.Sprintf("sessions[%d].end_date_time", i)
				errorMessages[key] = "end_date_time sessions must not below start_date_time"
			}
			if session.StartDateTime.Before(req.StartDate) {
				key := fmt.Sprintf("sessions[%d].start_date_time", i)
				errorMessages[key] = "start_date sessions must not below start_date event"
			}
			if session.EndDateTime.After(req.EndDate) {
				key := fmt.Sprintf("sessions[%d].end_date_time", i)
				errorMessages[key] = "end_date sessions must not above end_start event"
			}
		}

		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}
	event, err := h.eventUC.Create(req, userAuth, url)
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
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Create Event successfully", event))
}

func (h *EventHandler) GetAll(c *fiber.Ctx) error {
	var query request.EventPaginateRequest
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(http.StatusBadRequest, "invalid query parameters", err))
	}
	userAuth := c.Locals("user").(middleware.AuthUser)

	events, total, err := h.eventUC.GetAll(query, &userAuth.Tenant.ID)
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
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(fiber.StatusOK, "Get Event Successfully", events, query.Page, query.PageSize, total))
}

func (h *EventHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	event, err := h.eventUC.GetByID(id)
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
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Get Event Successfully", event))
}
