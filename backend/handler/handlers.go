package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"torospace.csudh.edu/api/sqlite"
)

var (
	db           *sqlite.DB
	sessionStore *session.Store
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env: %s", err)
	}

	var err error
	db, err = sqlite.NewDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	sessionStore = session.New(session.Config{
		Expiration:     30 * time.Minute,
		CookieHTTPOnly: true,
		// CookieSecure:  true, // HTTPS only
	})
}

func HelloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
