package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetTopicsHandler(c *fiber.Ctx) error {
	topics, err := db.GetTopics()
	if err != nil {
		log.Println("Failed to get topics in GetTopicsHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(topics)
}
