package handler

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/usecase"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventUC usecase.EventUsecase
}

func NewEventHandler(eventUC usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUC: eventUC}
}

func (h *EventHandler) GetTags(c *fiber.Ctx) error {
	eventTags, err := h.eventUC.GetTags()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse("Get Event Tags successfully", eventTags))
}

func (h *EventHandler) GetCategories(c *fiber.Ctx) error {
	eventCategories, err := h.eventUC.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse("Get Event Categories successfully", eventCategories))
}

func (h *EventHandler) Create(c *fiber.Ctx) error {
	var req request.CreateEventRequest

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

	if err := h.eventUC.Create(req, userAuth, method); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), err))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Create Event successfully"))
}
