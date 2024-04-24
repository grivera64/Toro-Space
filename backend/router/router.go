package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"torospace.csudh.edu/api/handler"
)

func SetupRoutes(app *fiber.App) {
	// Endpoint: /
	app.Get("/", handler.HelloHandler)
	app.Get("/monitor", monitor.New(monitor.Config{
		Title: "Toro Space Monitor",
	}))

	app.Get("/account/:accountID", handler.GetAccountHandler)
	app.Get("/account/:accountID/user/:userID", handler.GetUserHandler)
	app.Put("/account/:accountID/user/:userID/select", handler.SelectUserHandler)
	app.Post("/account/:accountID/user/:userID/post", handler.CreatePostHandler)

	app.Get("/user/self", handler.GetCurrentUserHandler)

	// Endpoint: /posts
	app.Get("/posts", handler.GetPostsHandler)
	app.Get("/posts/:postID", handler.GetPostHandler)
	app.Post("/posts/:postID/like", handler.LikePostHandler)

	// Endpoint: /topics
	app.Get("/topics", handler.GetTopicsHandler)

	// Endpoint: /auth/google
	app.Get("/auth/google", handler.GoogleAuthHandler)
	app.Get("/auth/google/callback", handler.GoogleAuthCallbackHandler)

	app.Get("/admin", handler.IsAdminHandler)
	app.Post("/admin/new/user", handler.CreateUserHandler)
	app.Get("/admin/account/:accountID", handler.GetAccountAdminHandler)
	app.Post("/admin/new/topic/:topicName", handler.CreateTopicHandler)

	app.Get("/logout", handler.LogoutHandler)

	// app.Use("*", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(fiber.StatusNotFound)
	// })
}
