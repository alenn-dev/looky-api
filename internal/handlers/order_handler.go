package handlers

import (
	"looky/internal/models"
	"looky/internal/services"

	"github.com/gofiber/fiber/v3"
)

func GetOrders(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := models.UserRole(c.Locals("role").(string))

	orders, err := services.GetOrders(userID, role)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(orders)
}

func GetOrder(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := models.UserRole(c.Locals("role").(string))

	order, err := services.GetOrder(id, userID, role)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(order)
}

func CreateOrder(c fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)

	var dto models.CreateOrderDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	order, err := services.CreateOrder(customerID, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func UpdateOrderStatus(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := models.UserRole(c.Locals("role").(string))

	var dto models.UpdateOrderStatusDTO
	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parameters"})
	}

	order, err := services.UpdateOrderStatus(id, userID, role, dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(order)
}
