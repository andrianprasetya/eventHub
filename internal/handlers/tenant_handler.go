package handlers

import (
	"github.com/andrianprasetya/eventHub/internal/usecases"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TenantHandler struct {
	TenantUC *usecases.TenantUsecase
}

func NewTenantHandler(tenantUC *usecases.TenantUsecase) *TenantHandler {
	return &TenantHandler{TenantUC: tenantUC}
}

func (h *TenantHandler) CreateTenant(c echo.Context) error {
	req := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Logo  string `json:"logo"`
		Plan  string `json:"plan"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	tenant, err := h.TenantUC.CreateTenant(req.Name, req.Email, req.Logo, req.Plan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tenant"})
	}

	return c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandler) GetTenantByID(c echo.Context) error {
	id := c.Param("id")
	tenant, err := h.TenantUC.GetTenantByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
	}
	return c.JSON(http.StatusOK, tenant)
}
