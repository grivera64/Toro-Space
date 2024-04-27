package handler

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/sqlite"
	"torospace.csudh.edu/api/util"
)

func GetPostsHandler(c *fiber.Ctx) error {
	postParams := &sqlite.PostParams{
		Before:      c.Query("before", ""),
		After:       c.Query("after", ""),
		PageSize:    c.QueryInt("pageSize", 10),
		SearchQuery: c.Query("search_query", ""),
	}

	posts, err := db.GetPosts(postParams)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(posts)
}

func GetPostHandler(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("postID")
	if err != nil {
		log.Println("Failed to get postID from params in GetPostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	post, err := db.GetPost(uint(postID))
	if err != nil {
		log.Println("Failed to get post by ID in GetPostHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(post)
}

func LikePostHandler(c *fiber.Ctx) error {

	like := c.Query("type", "like")

	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session store in LikePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in LikePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessUserID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in LikePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	postID, err := c.ParamsInt("postID")
	if err != nil {
		log.Println("Failed to get postID from params in LikePostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Println("Failed to get account by ID in LikePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: sessUserID})
	if err != nil {
		log.Println("Failed to get user by ID in LikePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if like == "like" {
		log.Println("Liking...")
		if err := db.AddLikeToPost(uint(postID), &user); err != nil {
			log.Println("Failed to like post in LikePostHandler")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else if like == "unlike" {
		log.Println("unliking...")
		if err := db.RemoveLikeFromPost(uint(postID), &user); err != nil {
			log.Println("Failed to unlike post in LikePostHandler", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else {
		log.Println("Invalid type in LikePostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	post, err := db.GetPost(uint(postID))
	if err != nil {
		log.Println("Failed to get post by ID in LikePostHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(post)
}

func CreatePostHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session store in CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
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
		log.Println("Failed to get userID from params in CreatePostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	sessUserID, ok := sess.Get("userID").(uint)
	if !ok || uint(userID) != sessUserID {
		log.Println("No userID in session for CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if role, ok := sess.Get("userRole").(entity.Role); !ok || (role != entity.RoleAdmin && role != entity.RoleOrganization) {
		log.Println("User is not an admin or organization in CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sess.SetExpiry(30 * time.Minute)
	if err := sess.Save(); err != nil {
		log.Println("Failed to save sesion store in CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	reqBody := fiber.Map{}
	if err := c.BodyParser(&reqBody); err != nil {
		log.Printf("Failed to parse request body in CreatePostHandle: %s", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	account, err := db.GetAccountByID(sessAccountID)
	if err != nil {
		log.Println("Failed to get account by ID in CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := util.BinarySearch(account.Users, entity.User{ID: sessUserID})
	if err != nil {
		log.Println("Failed to get user by ID in CreatePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if user.Role != entity.RoleAdmin && user.Role != entity.RoleOrganization {
		log.Println("User is not an admin or organization in CreatePostHandler (from DB)")
		return c.SendStatus(fiber.StatusForbidden)
	}

	postContent, ok := reqBody["content"].(string)
	if !ok || len(postContent) < 1 {
		log.Printf("Failed to get content from request body in CreatePostHandler: %v is %T", reqBody["content"], reqBody["content"])
		return c.SendStatus(fiber.StatusBadRequest)
	}

	post := &entity.Post{
		Content: strings.TrimSpace(postContent),
		Author:  user,
	}

	reqTopics, ok := reqBody["topics"].([]interface{})
	if ok {
		topics := []entity.Topic{}
		for _, topicName := range reqTopics {
			topicName, ok := topicName.(string)
			if !ok {
				log.Printf("Failed to get topic name from request body in CreatePostHandler: %v is %T", topicName, topicName)
				return c.SendStatus(fiber.StatusBadRequest)
			}
			if topic, err := db.GetTopicByName(topicName); err != nil {
				log.Printf("Failed to find %v in CreatePostHandler: %v", topicName, err)
			} else {
				topics = append(topics, *topic)
			}
		}
		post.Topics = topics
	} else {
		log.Println("Failed to get topics from request body in CreatePostHandler, ignoring...")
	}

	if err := db.AddPost(post); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(post)
}
