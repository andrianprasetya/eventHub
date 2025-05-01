package handler

import (
	"github.com/andrianprasetya/eventHub/internal/user"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUC user.UserUsecase
}

func NewUserHandler(userUC user.UserUsecase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {

	return nil
}
