package middleware

import (
	"looky/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func GraphQLMiddleware(c fiber.Ctx) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token != "" {
		claims, err := utils.ValidateJWT(token)
		if err == nil {
			c.Locals("user_id", claims.UserID)
			c.Locals("role", claims.Role)
		}
	}
	return c.Next()
}
