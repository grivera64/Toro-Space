package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/gateway/googleoauth"
	"torospace.csudh.edu/api/mapper"
	"torospace.csudh.edu/api/util"
)

var (
	googleGateway googleoauth.GoogleOauthGateway
)

func init() {
	googleGateway = googleoauth.NewV2()
}

func GoogleAuthHandler(c *fiber.Ctx) error {
	// Check if user is already authenticated
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	accountID, ok := session.Get("accountID").(uint)
	if ok {
		// Redirect to the login page (or wherever you want to send the user after they log out)
		return c.Redirect(fmt.Sprintf("http://localhost:3030/account/%d", accountID))
		// return c.Redirect("http://localhost:3000/select")
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
	session.SetExpiry(30 * time.Minute)
	if err := session.Regenerate(); err != nil {
		log.Printf("Failed to regenerate session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	session.Set("accountID", account.ID)

	if account.ID != session.Get("accountID") {
		log.Println("User ID is not being set properly (sess not set properly)")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if account.ID == 0 {
		log.Println("User ID is not being set properly (0)")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Redirect to the login page (or wherever you want to send the user after they log out)
	return c.Redirect(fmt.Sprintf("http://localhost:3030/account/%d", account.ID))
	// return c.Redirect("http://localhost:3000/select")
}

func SelectUserHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	sess.SetExpiry(30 * time.Minute)

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Printf("Failed to get userID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	accountID, err := c.ParamsInt("accountID")
	if err != nil {
		if c.Params("accountID") != "self" {
			log.Println("Failed to get accountID from params")
			return c.SendStatus(fiber.StatusBadRequest)
		}
	} else if uint(accountID) != sessAccountID {
		log.Printf("Failed to get accountID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID, err := c.ParamsInt("userID")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: uint(userID)})
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	sess.Set("userID", user.ID)
	sess.Set("userRole", user.Role)
	if err := sess.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
}

func GetCurrentUserHandler(c *fiber.Ctx) error {
	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	sessAccountID, ok := session.Get("accountID").(uint)
	if !ok {
		log.Printf("Failed to get accountID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	sessUserID, ok := session.Get("userID").(uint)
	if !ok {
		log.Printf("Failed to get userID from session")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: sessUserID})
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
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
	return c.Redirect("http://localhost:3030/")
	// return c.Redirect("http://localhost:3000")
}
