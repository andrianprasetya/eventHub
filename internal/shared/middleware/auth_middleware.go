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
	Tenant   TenantPayload `json:"tenant"`
	Role     RolePayload   `json:"role"`
	Token    string        `json:"token"`
}

type RolePayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type TenantPayload struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	LogoUrl  string `json:"logo_url"`
	Domain   string `json:"domain"`
	IsActive int    `json:"is_active"`
}

func AuthMiddleware(redis redisser.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(fiber.StatusUnauthorized, "Missing or invalid token", nil))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check token in Redis
		ctx := context.Background()
		data, err := redis.Get(ctx, "user:jwt:"+tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(fiber.StatusUnauthorized, "Token expired or invalid", err))
		}

		var authUser AuthUser
		if err := json.Unmarshal([]byte(data), &authUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, "Something went wrong", err))
		}
		// save payload ke context
		c.Locals("user", authUser)
		// Token is valid
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authUser := c.Locals("user") // assuming you already store the user info in Locals after auth

		if authUser == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(fiber.StatusUnauthorized, "Unauthorized", nil))
		}

		// assume user struct has a Role field like: authUser.(User).Role
		user := authUser.(AuthUser) // You can define your own struct or DTO

		for _, role := range roles {
			if role == user.Role.Slug {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse(fiber.StatusForbidden, "Forbidden: insufficient permissions", nil))
	}
}
