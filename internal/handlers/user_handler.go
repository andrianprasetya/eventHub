package handlers

import (
	"github.com/andrianprasetya/eventHub/internal/usecases"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	UserUC *usecases.UserUsecase
}

func NewUserHandler(userUC *usecases.UserUsecase) *UserHandler {
	return &UserHandler{UserUC: userUC}
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	req := struct {
		TenantID string `json:"tenant_id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.UserUC.RegisterUser(req.TenantID, req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.UserUC.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}
