package middleware

import (
	"context"
	"encoding/json"
	"github.com/andrianprasetya/eventHub/internal/shared/redisser"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthUser struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	IsActive int           `json:"is_active"`
	Tenant   RolePayload   `json:"tenant"`
	Role     TenantPayload `json:"role"`
	Token    string        `json:"token"`
}

type RolePayload struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Slug        string `json:"Slug"`
	Description string `json:"Description"`
	CreatedAt   string `json:"CreatedAt"`
	UpdatedAt   string `json:"UpdatedAt"`
}

type TenantPayload struct {
	ID        string `json:"ID"`
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	LogoUrl   string `json:"LogoUrl"`
	Domain    string `json:"Domain"`
	IsActive  int    `json:"IsActive"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

func AuthMiddleware(redis redisser.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse("Missing or invalid token", nil))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check token in Redis
		ctx := context.Background()
		data, err := redis.Get(ctx, "user:jwt:"+tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse("Token expired or invalid", err))
		}

		var authUser AuthUser
		if err := json.Unmarshal([]byte(data), &authUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Something went wrong", err))
		}
		// save payload ke context
		c.Locals("user", authUser)
		// Token is valid
		return c.Next()
	}
}
