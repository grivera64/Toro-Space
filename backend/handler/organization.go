package handler

import (
	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/sqlite"
)

func GetOrganizationsHandler(c *fiber.Ctx) error {
	organizationParams := &sqlite.OrganizationParams{
		Before:      c.Query("before", ""),
		After:       c.Query("after", ""),
		PageSize:    c.QueryInt("page_size", 10),
		SearchQuery: c.Query("search_query", ""),
	}
	organizationsResult, err := db.GetOrganizations(organizationParams)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(organizationsResult)
}

func GetOrganizationHandler(c *fiber.Ctx) error {
	organizationID, err := c.ParamsInt("organizationID")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	organization, err := db.GetOrganization(uint(organizationID))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(organization)
}
