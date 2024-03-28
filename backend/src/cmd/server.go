package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"torospace.csudh.edu/api/src/router"
)

func main() {
	// Fiber Setup
	app := fiber.New(fiber.Config{})

	// Add healthcheck middleware for /livez and /readyz
	app.Use(healthcheck.New(healthcheck.Config{}))

	// Add routes
	router.SetupRoutes(app)

	if err := app.Listen(":3030"); err != nil {
		log.Fatal(err)
	}
}
