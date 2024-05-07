package handler

import (
	"context"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/entity"
	pb "torospace.csudh.edu/api/proto/spam_detector"
	"torospace.csudh.edu/api/sqlite"
	"torospace.csudh.edu/api/util"
)

func GetPostsHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		sess = nil
	}

	var userRole any
	var ok bool
	if sess != nil {
		userRole, ok = sess.Get("userRole").(entity.Role)
	}
	postParams := &sqlite.PostParams{
		Before:      c.Query("before", ""),
		After:       c.Query("after", ""),
		PageSize:    c.QueryInt("page_size", 10),
		SearchQuery: c.Query("search_query", ""),
		GetHidden:   sess != nil && ((ok && userRole == entity.RoleAdmin) || (userRole == entity.RoleOrganization)),
	}

	postsResult, err := db.GetPosts(postParams)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(postsResult)
}

func GetPostsByOrganizationHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		sess = nil
	}

	organizationID, err := c.ParamsInt("organizationID")
	if err != nil {
		log.Println("Failed to get organizationID from params in GetPostsByOrganizationHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var userRole any
	var user *entity.User
	var ok bool
	if sess != nil {
		userRole, ok = sess.Get("userRole").(entity.Role)
		if ok {
			user, ok = sess.Get("userID").(*entity.User)
		}
	}
	postParams := &sqlite.PostParams{
		Before:      c.Query("before", ""),
		After:       c.Query("after", ""),
		PageSize:    c.QueryInt("page_size", 10),
		SearchQuery: c.Query("search_query", ""),
		GetHidden:   sess != nil && ((ok && userRole == entity.RoleAdmin) || (user.ID == uint(organizationID))),
	}

	postsResult, err := db.GetPostsByOrganization(uint(organizationID), postParams)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(postsResult)
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

func DeletePostHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session store in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessUserID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
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

	sess.SetExpiry(30 * time.Minute)
	if err := sess.Save(); err != nil {
		log.Println("Failed to save sesion store in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	postID, err := c.ParamsInt("postID")
	if err != nil {
		log.Println("Failed to get postID from params in DeletePostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	post, err := db.GetPost(uint(postID))
	if err != nil {
		log.Println("Failed to get post by ID in DeletePostHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user.Role != entity.RoleAdmin && post.Author.ID != sessUserID {
		log.Println("User is not the author of the post in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := db.DeletePost(uint(postID)); err != nil {
		log.Println("Failed to delete post in DeletePostHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func HidePostHandler(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Println("Failed to get session store in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessAccountID, ok := sess.Get("accountID").(uint)
	if !ok {
		log.Println("Failed to get accountID in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	sessUserID, ok := sess.Get("userID").(uint)
	if !ok {
		log.Println("Failed to get userID in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
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

	sess.SetExpiry(30 * time.Minute)
	if err := sess.Save(); err != nil {
		log.Println("Failed to save sesion store in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	postID, err := c.ParamsInt("postID")
	if err != nil {
		log.Println("Failed to get postID from params in DeletePostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	post, err := db.GetPost(uint(postID))
	if err != nil {
		log.Println("Failed to get post by ID in DeletePostHandler")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user.Role != entity.RoleAdmin && post.Author.ID != sessUserID {
		log.Println("User is not the author of the post in DeletePostHandler")
		return c.SendStatus(fiber.StatusForbidden)
	}

	action := c.Query("action", "hide")
	if action == "hide" {
		if err := db.HidePost(uint(postID)); err != nil {
			log.Println("Failed to delete post in DeletePostHandler")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else if action == "unhide" {
		if err := db.UnhidePost(uint(postID)); err != nil {
			log.Println("Failed to delete post in DeletePostHandler")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else {
		log.Println("Invalid action in HidePostHandler")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
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
	postContent = strings.TrimSpace(postContent)

	conn, err := grpc.NewClient("127.0.0.1:3060", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		client := pb.NewSpamDetectorClient(conn)
		detectorResp, err := client.Scan(context.Background(), &pb.ScanRequest{Content: postContent})
		log.Printf("Detector says: %s", detectorResp.Result.String())
		if err == nil {
			if detectorResp.Result == pb.ScanResponse_SPAM {
				return c.SendStatus(fiber.StatusNotAcceptable)
			}
		}
		conn.Close()
	} else {
		log.Printf("Failed to connect to grpc: %v", err)
	}

	post := &entity.Post{
		Content: postContent,
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
