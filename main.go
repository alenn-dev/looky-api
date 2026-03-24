package main

import (
	"log"
	"looky/internal/config"
	"looky/internal/database"
	"looky/internal/handlers"
	"looky/internal/kafka"
	"looky/internal/middleware"
	"looky/internal/routes"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnv()
	database.Connect()

	brokers := strings.Split(config.Env.KafkaBrokers, ",")
	if err := kafka.InitProducer(brokers); err != nil {
		log.Fatal("couldn't init kafka producer:", err)
	}
	defer kafka.CloseProducer()
	kafka.StartConsumer(brokers)

	app := fiber.New()

	app.Get("/ws", handlers.WebSocketHandler)
	app.Use("/graphql", middleware.GraphQLMiddleware)
	app.Post("/graphql", handlers.GraphQLHandler())
	// app.Get("/playground", handlers.GraphQLPlayground) // only devs

	routes.Setup(app)

	log.Fatal(app.Listen(":3000"))
}
