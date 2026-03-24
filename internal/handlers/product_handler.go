package handlers

import (
	"looky/internal/models"
	"looky/internal/services"

	"github.com/gofiber/fiber/v3"
)

func GetProducts(c fiber.Ctx) error {
	restaurantID := c.Params("restaurantId")
	products, err := services.GetProducts(restaurantID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(products)
}

func GetProduct(c fiber.Ctx) error {
	id := c.Params("id")
	restaurantID := c.Params("restaurantId")
	product, err := services.GetProduct(id, restaurantID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(product)
}

func CreateProduct(c fiber.Ctx) error {
	restaurantID := c.Params("restaurantId")
	ownerID := c.Locals("user_id").(string)

	var dto models.CreateProductDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	product, err := services.CreateProduct(restaurantID, ownerID, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")
	restaurantID := c.Params("restaurantId")
	ownerID := c.Locals("user_id").(string)

	var dto models.UpdateProductDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	product, err := services.UpdateProduct(id, restaurantID, ownerID, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(product)
}

func DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")
	restaurantID := c.Params("restaurantId")
	ownerID := c.Locals("user_id").(string)

	if err := services.DeleteProduct(id, restaurantID, ownerID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
