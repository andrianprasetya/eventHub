package handler

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/usecase"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type EventHandler struct {
	eventUC usecase.EventUsecase
}

func NewEventHandler(eventUC usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUC: eventUC}
}

func (h *EventHandler) GetTags(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	eventTags, total, err := h.eventUC.GetTags(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(fiber.StatusOK, "Get Event Tags successfully", eventTags, page, pageSize, total))
}

func (h *EventHandler) GetCategories(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	eventCategories, total, err := h.eventUC.GetCategories(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(fiber.StatusOK, "Get Event Categories successfully", eventCategories, page, pageSize, total))
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationResponse(fiber.StatusBadRequest, errorMessages))
	}
	event, err := h.eventUC.Create(req, userAuth, url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Create Event successfully", event))
}

func (h *EventHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	events, total, err := h.eventUC.GetAll(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(fiber.StatusOK, "Get Event Successfully", events, page, pageSize, total))
}

func (h *EventHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	event, err := h.eventUC.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "Get Event Successfully", event))
}
