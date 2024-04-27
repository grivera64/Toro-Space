package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/util"
)

func GetAccountHandler(c *fiber.Ctx) error {
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

	session.SetExpiry(30 * time.Minute)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Printf("Failed to get user from database: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.JSON(account.Users)
}

func GetUserHandler(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("userID", -1)
	if err != nil || userID < 1 {
		log.Printf("Failed to get userID from params: %s", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	session, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	session.SetExpiry(30 * time.Minute)

	sessAccountID, ok := session.Get("accountID").(uint)
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

	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
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
	return c.JSON(user)
}
