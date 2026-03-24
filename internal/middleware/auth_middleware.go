package middleware

import (
	"looky/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(400).JSON(fiber.Map{"error": "unauthorized"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)
	return c.Next()
}
