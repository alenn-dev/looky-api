package handlers

import (
	"looky/internal/models"
	"looky/internal/services"

	"github.com/gofiber/fiber/v3"
)

func GetMe(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user id required"})
	}

	user, err := services.GetMe(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(user)
}

func UpdateUser(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user id required"})
	}

	var dto = models.UpdateUserDTO{}
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	user, err := services.UpdateUser(userID, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(user)
}
