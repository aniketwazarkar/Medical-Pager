package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"medical-pager/internal/api"
	"medical-pager/internal/db"
	"medical-pager/internal/redis"
	"medical-pager/utils"
)

func main() {
	// Initialize Config
	utils.LoadConfig()

	// Initialize Infrastructure
	db.Connect()
	db.SetupIndexes()
	redis.Connect()

	// Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Medical Pager v1.0",
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} \n\tReq: ${body}\n\tRes: ${resBody}\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust in production
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health Check Route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Medical Pager API is running",
		})
	})

	// Setup API Routes
	api.SetupRoutes(app)

	// Start Server
	port := utils.GetEnv("PORT", "5000")
	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
