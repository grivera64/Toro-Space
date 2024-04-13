package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/mapper"
	"torospace.csudh.edu/api/sqlite"
	"torospace.csudh.edu/api/util"
)

var (
	db           *sqlite.DB
	sessionStore *session.Store
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
		// return c.Redirect(fmt.Sprintf("http://localhost:3030/account/%d", accountID))
		return c.Redirect("http://localhost:3000/select")
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
	// return c.Redirect(fmt.Sprintf("http://localhost:3030/account/%d", account.ID))
	return c.Redirect("http://localhost:3000/select")
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
	// return c.Redirect("http://localhost:3030/")
	return c.Redirect("http://localhost:3000")
}

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

func GetPostsHandler(c *fiber.Ctx) error {
	postParams := &sqlite.PostParams{
		Before:   c.Query("before", ""),
		After:    c.Query("after", ""),
		PageSize: c.QueryInt("pageSize", 10),
	}

	posts, err := db.GetPosts(postParams)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(posts)
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

	reqTopicsInterface, ok := reqBody["topics"]
	var reqTopics []interface{}
	if ok {
		reqTopics, ok = reqTopicsInterface.([]interface{})
		if !ok {
			log.Printf("Failed to get topics from request body in CreatePostHandler: %v is %T", reqTopicsInterface, reqTopicsInterface)
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	topics := []entity.Topic{}
	for _, topicName := range reqTopics {
		topicName, ok := topicName.(string)
		if !ok {
			log.Printf("Failed to get topic name from request body in CreatePostHandler: %v is %T", topicName, topicName)
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if topic, err := db.GetTopicByName(topicName); err != nil {
			topics = append(topics, *topic)
		} else {
			log.Printf("Failed to find %v in CreatePostHandler", topicName)
		}
	}

	post := &entity.Post{
		Content: reqBody["content"].(string),
		Author:  user,
	}

	if len(topics) > 0 {
		post.Topics = topics
	}

	if err := db.AddPost(post); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(post)
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
