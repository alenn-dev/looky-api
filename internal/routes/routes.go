package routes

import (
	"looky/internal/handlers"
	"looky/internal/middleware"
	"looky/internal/models"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	api.Use(middleware.RateLimiter())

	auth := api.Group("/auth")
	auth.Use(middleware.AuthRateLimiter())

	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	user := api.Group("/user", middleware.AuthMiddleware)

	user.Get("/", handlers.GetMe)
	user.Put("/", handlers.UpdateUser)

	restaurants := api.Group("/restaurants", middleware.AuthMiddleware)

	restaurants.Get("/", handlers.GetRestaurants)
	restaurants.Get("/:id", handlers.GetRestaurant)
	restaurants.Post("/", middleware.RequireRole(models.RestaurantOwner), handlers.CreateRestaurant)
	restaurants.Patch("/:id", middleware.RequireRole(models.RestaurantOwner), handlers.UpdateRestaurant)
	restaurants.Delete("/:id", middleware.RequireRole(models.RestaurantOwner), handlers.DeleteRestaurant)

	products := api.Group("/restaurants/:restaurantId/products", middleware.AuthMiddleware)

	products.Get("/", handlers.GetProducts)
	products.Get("/:id", handlers.GetProduct)
	products.Post("/", middleware.RequireRole(models.RestaurantOwner), handlers.CreateProduct)
	products.Patch("/:id", middleware.RequireRole(models.RestaurantOwner), handlers.UpdateProduct)
	products.Delete("/:id", middleware.RequireRole(models.RestaurantOwner), handlers.DeleteProduct)

	orders := api.Group("/orders", middleware.AuthMiddleware)

	orders.Get("/", handlers.GetOrders)
	orders.Get("/:id", handlers.GetOrder)
	orders.Post("/", middleware.RequireRole(models.Customer), handlers.CreateOrder)
	orders.Patch("/:id/status", handlers.UpdateOrderStatus)
}
