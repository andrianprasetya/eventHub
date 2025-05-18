package handler

import (
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type RoleHandler struct {
	roleUC usecase.RoleUsecase
}

func NewRoleHandler(roleUC usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUC: roleUC}
}

func (u *RoleHandler) GetAll(c *fiber.Ctx) error {
	var query request.RolePaginateParams

	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(http.StatusBadRequest, "invalid query parameters", err))
	}

	roles, total, err := u.roleUC.GetAll(query)
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
		"Get Role successfully",
		roles,
		query.Page,
		query.PageSize,
		total,
	))
}

func (u *RoleHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := u.roleUC.GetByID(id)
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "get Role successfully", role))
}
