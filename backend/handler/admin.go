package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/util"
)

func IsAdminHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session store in IsAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in IsAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessUserID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in IsAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Println("Failed to get account by ID in IsAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: sessUserID})
	if err != nil {
		log.Println("Failed to get user by ID in IsAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if user.Role != entity.RoleAdmin {
		log.Println("User is not an admin in IsAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.SendStatus(fiber.StatusOK)
}

func CreateTopicHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session store in CreateTopicHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in CreateTopicHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}
	sessUserID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in CreateTopicHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if role, ok := sess.Get("userRole").(entity.Role); !ok || role != entity.RoleAdmin {
		log.Println("User is not an admin in CreateTopicHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Println("Failed to get account by ID in CreateTopicHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: sessUserID})
	if err != nil {
		log.Println("Failed to get user by ID in CreateTopicHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if user.Role != entity.RoleAdmin {
		log.Println("User is not an admin in CreateTopicHandler (from DB)")
		return c.SendStatus(fiber.StatusForbidden)
	}

	topicName := c.Params("topicName")
	if len(topicName) < 1 {
		log.Println("Failed to get topicName from params in CreateTopicHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	topic := &entity.Topic{
		Name: topicName,
	}

	if err := db.CreateTopic(topic); err != nil {
		log.Println("Failed to create topic in CreateTopicHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(topic)
}

func CreateUserHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session in CreateUserHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}
	accountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in CreateUserHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}
	userID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in CreateUserHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	account, err := db.GetAccountByID(accountID)
	if err != nil {
		log.Println("Failed to get account by ID in CreateUserHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: userID})
	if err != nil {
		log.Println("Failed to get user by ID in CreateUserHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessRole, ok := sess.Get("userRole").(entity.Role)
	if !ok || (user.Role != entity.RoleAdmin || sessRole != entity.RoleAdmin) {
		log.Printf("User is not an admin in CreateUserHandler (user.Role = %s, sessRole = %s)", user.Role, sessRole)
		return c.SendStatus(fiber.StatusForbidden)
	}

	reqBody := struct {
		AccountID   uint        `json:"account_id"`
		DisplayName string      `json:"display_name"`
		AvatarUrl   string      `json:"avatar_url"`
		Role        entity.Role `json:"role"`
	}{}
	if err := c.BodyParser(&reqBody); err != nil {
		log.Printf("Failed to parse request body in CreateUserHandler: %s", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	newUserAccount, err := db.GetAccountByID(reqBody.AccountID)
	if err != nil {
		log.Println("Failed to get account by ID in CreateUserHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}
	newUser := &entity.User{
		DisplayName: reqBody.DisplayName,
		AvatarUrl:   reqBody.AvatarUrl,
		Role:        reqBody.Role,
	}
	if err := db.AddAccountUser(newUserAccount, newUser); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"message": "User created successfully",
			"user":    newUser,
		})
}

func GetAccountAdminHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session in GetAccountAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	getAccountID, err := c.ParamsInt("accountID", -1)
	if err != nil || getAccountID < 1 {
		log.Println("accountID was not parsable in params, or is negative")
		return c.SendStatus(fiber.StatusBadRequest)
	}
	accountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in GetAccountAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}
	userID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in GetAccountAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	account, err := db.GetAccountByID(accountID)
	if err != nil {
		log.Println("Failed to get current account by ID in GetAccountAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: userID})
	if err != nil {
		log.Println("Failed to get current user by ID in GetAccountAdminHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessRole, ok := sess.Get("userRole").(entity.Role)
	if !ok || (user.Role != entity.RoleAdmin || sessRole != entity.RoleAdmin) {
		log.Printf("User is not an admin in GetAccountAdminHandler (user.Role = %s, sessRole = %s)", user.Role, sessRole)
		return c.SendStatus(fiber.StatusForbidden)
	}

	getUserAccount, err := db.GetAccountByID(uint(getAccountID))
	if err != nil {
		log.Printf("Failed to get requested account by ID in GetAccountAdminHandler (accountID = %d)", getAccountID)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(fiber.Map{
		"account": getUserAccount,
	})
}
