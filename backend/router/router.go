package router

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"torospace.csudh.edu/api/handler"
)

var (
	limitedCtx    = map[string]context.Context{}
	limitedCancel = map[string]context.CancelFunc{}
)

func limiterNext(c *fiber.Ctx) bool {
	if ctx, ok := limitedCtx[c.IP()]; ok {
		select {
		case <-ctx.Done():
			limitedCancel[c.IP()]()
			delete(limitedCancel, c.IP())
			delete(limitedCtx, c.IP())
			return true
		default:
			return false
		}
	}
	return false
}

func limiterReached(c *fiber.Ctx) error {
	if _, ok := limitedCtx[c.IP()]; ok {
		return c.SendStatus(fiber.StatusTooManyRequests)
	}
	limitedCtx[c.IP()], limitedCancel[c.IP()] = context.WithTimeout(context.Background(), 1*time.Hour)
	return c.SendStatus(fiber.StatusTooManyRequests)
}

func SetupRoutes(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:          3000,
		Expiration:   1 * time.Hour,
		Next:         limiterNext,
		LimitReached: limiterReached,
	}))

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
	app.Delete("/posts/:postID", handler.DeletePostHandler)
	app.Put("/posts/:postID", handler.HidePostHandler)
	app.Post("/posts/:postID/like", handler.LikePostHandler)

	// Endpoint: /topics
	app.Get("/topics", handler.GetTopicsHandler)

	// Endpoint: /organizations
	app.Get("/organizations", handler.GetOrganizationsHandler)
	app.Get("/organizations/:organizationID", handler.GetOrganizationHandler)
	app.Get("/organizations/:organizationID/posts", handler.GetPostsByOrganizationHandler)

	// Endpoint: /auth/google
	app.Get("/auth/google", handler.GoogleAuthHandler)
	app.Get("/auth/google/callback", handler.GoogleAuthCallbackHandler)

	// Endpoint: /logout
	app.Get("/logout", handler.LogoutHandler)

	app.Get("/admin", handler.IsAdminHandler)
	app.Post("/admin/new/user", handler.CreateUserHandler)
	app.Get("/admin/account/:accountID", handler.GetAccountAdminHandler)
	app.Post("/admin/new/topic/:topicName", handler.CreateTopicHandler)

	// app.Use("*", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(fiber.StatusNotFound)
	// })
}
