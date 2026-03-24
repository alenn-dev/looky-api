package middleware

import (
	"looky/internal/models"

	"github.com/gofiber/fiber/v3"
)

func RequireRole(roles ...models.UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := models.UserRole(c.Locals("role").(string))
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
}
