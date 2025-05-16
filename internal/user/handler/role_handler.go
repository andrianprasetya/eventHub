package handler

import (
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type RoleHandler struct {
	roleUC usecase.RoleUsecase
}

func NewRoleHandler(roleUC usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUC: roleUC}
}

func (u *RoleHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	roles, total, err := u.roleUC.GetAll(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateDataResponse(
		fiber.StatusOK,
		"Get Role successfully",
		roles,
		page,
		pageSize,
		total,
	))
}

func (u *RoleHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := u.roleUC.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, err.Error(), err))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDataResponse(fiber.StatusOK, "get Role successfully", role))
}
