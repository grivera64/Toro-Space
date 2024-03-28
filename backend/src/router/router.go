package router

import (
	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/src/handler"
)

func SetupRoutes(app *fiber.App) {
	// Endpoint: /
	app.Get("/", handler.HelloHandler)
}
