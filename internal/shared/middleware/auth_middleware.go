package middleware

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/shared/redisser"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware(redis redisser.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check token in Redis
		ctx := context.Background()
		_, err := redis.Get(ctx, "jwt:"+tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired or invalid",
			})
		}

		// Token is valid
		return c.Next()
	}
}
