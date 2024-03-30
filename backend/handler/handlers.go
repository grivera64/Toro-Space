package handler

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/mapper"
	"torospace.csudh.edu/api/sqlite"
)

var (
	oauthConfig  *oauth2.Config
	db           *sqlite.DB
	sessionStore *session.Store
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env: %s", err)
	}
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("G_CLIENT_ID"),
		ClientSecret: os.Getenv("G_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("G_REDIRECT"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
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

func GoogleAuthHandler(c *fiber.Ctx) error {
	// TODO: Implement a secure Auth state challenge
	googleUrl := oauthConfig.AuthCodeURL("not-implemented-yet")
	return c.Redirect(googleUrl)
}

func GoogleAuthCallbackHandler(c *fiber.Ctx) error {
	// Get Auth Code (?code=...)
	code := c.Query("code")

	// Exchange auth code for token
	token, err := oauthConfig.Exchange(c.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange auth code for token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Map to a User
	googleAuth, err := mapper.TokenToGoogleAuth(token.AccessToken)
	if err != nil {
		log.Printf("Failed to map token to GoogleAuth: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var user *entity.User
	user, err = db.GetUserByGoogleID(googleAuth.ID)

	// User does not exist, create a new user
	if err != nil {
		user = mapper.GoogleAuthToUser(googleAuth)
		if err := db.AddUser(user); err != nil {
			log.Printf("Failed to add user to the database: %s", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Store user's internal ID in session
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	session.Set("userID", user.ID)
	if err := session.Regenerate(); err != nil {
		log.Printf("Failed to regenerate session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// TODO: Set to redirect to frontend, temporary redirect to /user
	return c.Redirect("http://localhost:3000")
}

func LogoutHandler(c *fiber.Ctx) error {
	// Get the current session
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Destroy the session
	if err := session.Destroy(); err != nil {
		log.Printf("Failed to destroy session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Redirect to the login page (or wherever you want to send the user after they log out)
	return c.Redirect("http://localhost:3000")
}

func GetUserHandler(c *fiber.Ctx) error {
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID, ok := session.Get("userID").(uint)
	if !ok {
		log.Printf("Failed to get userID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	session.SetExpiry(30 * time.Minute)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
}
