package handlers

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"github.com/andrianprasetya/eventHub/internal/usecases"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Register(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	savedUser, err := usecases.RegisterUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, savedUser)
}

func Login(c echo.Context) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credentials"})
	}

	token, err := usecases.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
