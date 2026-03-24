package handlers

import (
	"looky/internal/models"
	"looky/internal/services"

	"github.com/gofiber/fiber/v3"
)

func GetRestaurants(c fiber.Ctx) error {
	restaurants, err := services.GetRestaurants()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(restaurants)
}

func GetRestaurant(c fiber.Ctx) error {
	id := c.Params("id")
	restaurant, err := services.GetRestaurant(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(restaurant)
}

func CreateRestaurant(c fiber.Ctx) error {
	ownerID := c.Locals("user_id").(string)

	var dto models.CreateRestaurantDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	restaurant, err := services.CreateRestaurant(ownerID, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(restaurant)
}

func UpdateRestaurant(c fiber.Ctx) error {
	id := c.Params("id")
	ownerID := c.Locals("user_id").(string)

	var dto models.UpdateRestaurantDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	restaurant, err := services.UpdateRestaurant(id, ownerID, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(restaurant)
}

func DeleteRestaurant(c fiber.Ctx) error {
	id := c.Params("id")
	ownerID := c.Locals("user_id").(string)

	if err := services.DeleteRestaurant(id, ownerID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
