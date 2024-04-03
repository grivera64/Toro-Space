package router

import (
	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/handler"
)

func SetupRoutes(app *fiber.App) {
	// Endpoint: /
	app.Get("/", handler.HelloHandler)

	app.Get("/account/:accountID", handler.GetAccountHandler)
	app.Get("/account/:accountID/user/:userID", handler.GetUserHandler)
	app.Put("/account/:accountID/user/:userID/select", handler.SelectUserHandler)
	app.Post("/account/:accountID/user/:userID/post", handler.CreatePostHandler)

	// Endpoint: /posts
	app.Get("/posts", handler.GetPostsHandler)

	// Endpoint: /auth/google
	app.Get("/auth/google", handler.GoogleAuthHandler)
	app.Get("/auth/google/callback", handler.GoogleAuthCallbackHandler)

	app.Get("/logout", handler.LogoutHandler)

	// app.Use("*", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(fiber.StatusNotFound)
	// })
}
