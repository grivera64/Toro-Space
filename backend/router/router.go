package router

import (
	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/handler"
)

func SetupRoutes(app *fiber.App) {
	// Endpoint: /
	app.Get("/", handler.HelloHandler)

	app.Get("/user", handler.GetUserHandler)

	// Endpoint: /auth/google
	app.Get("/auth/google", handler.GoogleAuthHandler)
	app.Get("/auth/google/callback", handler.GoogleAuthCallbackHandler)
}
