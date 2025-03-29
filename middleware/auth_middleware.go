package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// JWTMiddleware untuk validasi token JWT
func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		TokenLookup: "header:Authorization",
		SigningKey:  []byte("secret"),
	})
}

// RoleMiddleware untuk validasi peran pengguna
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("user_role") // Role diambil dari JWT Claims

			if userRole == nil {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
			}

			roleStr := userRole.(string)
			for _, role := range allowedRoles {
				if strings.EqualFold(roleStr, role) {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized role"})
		}
	}
}
