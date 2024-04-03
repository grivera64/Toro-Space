package router

import (
	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/handler"
)

func SetupRoutes(app *fiber.App) {
	// Endpoint: /
	app.Get("/", handler.HelloHandler)

	app.Get("/account/:id", handler.GetAccountHandler)
	app.Get("/users/:id", handler.GetUserHandler)

	// Endpoint: /auth/google
	app.Get("/auth/google", handler.GoogleAuthHandler)
	app.Get("/auth/google/callback", handler.GoogleAuthCallbackHandler)

	app.Get("/logout", handler.LogoutHandler)
}
