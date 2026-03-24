package handlers

import (
	"looky/internal/models"
	"looky/internal/services"

	"github.com/gofiber/fiber/v3"
)

func Register(c fiber.Ctx) error {
	var dto models.RegisterDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	user, err := services.Register(dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(user)
}

func Login(c fiber.Ctx) error {
	var dto models.LoginDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	token, err := services.Login(dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(token)
}
