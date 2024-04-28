package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/sqlite"
)

func GetTopicsHandler(c *fiber.Ctx) error {
	topicParams := &sqlite.TopicParams{
		PageSize:    c.QueryInt("page_size", 10),
		SearchQuery: c.Query("search_query", ""),
	}
	topicsResult, err := db.GetTopics(topicParams)
	if err != nil {
		log.Println("Failed to get topics in GetTopicsHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(topicsResult)
}
