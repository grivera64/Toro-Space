package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/gateway/googleoauth"
	"torospace.csudh.edu/api/mapper"
	"torospace.csudh.edu/api/sqlite"
	"torospace.csudh.edu/api/util"
)

var (
	oauthConfig   *oauth2.Config
	db            *sqlite.DB
	sessionStore  *session.Store
	googleGateway googleoauth.GoogleOauthGateway
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

	googleGateway = googleoauth.NewV2()
}

func HelloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func GoogleAuthHandler(c *fiber.Ctx) error {
	// Check if user is already authenticated
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	_, ok := session.Get("accountID").(uint)
	if ok {
		// Redirect to the login page (or wherever you want to send the user after they log out)
		return c.Redirect("http://localhost:3030/users/0")
		// return c.Redirect("http://localhost:3000")
	}

	return c.Redirect(googleGateway.GetAuthUrl())
}

func GoogleAuthCallbackHandler(c *fiber.Ctx) error {
	// Get Auth Code (?code=...)
	code := c.Query("code")

	// Exchange auth code for token
	token, err := googleGateway.GetToken(c.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange auth code for token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Map to a User
	googleUser, err := googleGateway.GetUserInfo(token.AccessToken)
	if err != nil {
		log.Printf("Failed to map token to GoogleAuth: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var account *entity.Account
	account, err = db.GetAccountByGoogleID(googleUser.ID)

	// User does not exist, create a new user
	if err != nil {
		account = mapper.GoogleUserToAccount(googleUser)
		if err := db.AddAccount(account); err != nil {
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
	session.Set("accountID", account.ID)
	if err := session.Regenerate(); err != nil {
		log.Printf("Failed to regenerate session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Redirect to the login page (or wherever you want to send the user after they log out)
	return c.Redirect(fmt.Sprintf("http://localhost:3030/account/%d", account.ID))
	// return c.Redirect("http://localhost:3000")
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
	return c.Redirect("http://localhost:3030/users/0")
	// return c.Redirect("http://localhost:3000")
}

func GetAccountHandler(c *fiber.Ctx) error {
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	accountID, ok := session.Get("accountID").(uint)
	if !ok {
		log.Printf("Failed to get userID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	account, err := db.GetAccountByID(accountID)
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	session.SetExpiry(30 * time.Minute)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(account.Users)
}

func GetUserHandler(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Failed to get userID from params: %s", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	accountID, ok := session.Get("accountID").(uint)
	if !ok {
		log.Printf("Failed to get userID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	account, err := db.GetAccountByID(accountID)
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: uint(userID)})
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
