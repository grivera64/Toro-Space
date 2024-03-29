package mapper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"torospace.csudh.edu/api/entity"
)

func TokenToGoogleAuth(token string) (*entity.GoogleAuthPayload, error) {
	client := fiber.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", token))
	if client == nil {
		return nil, fmt.Errorf("unable to get user info")
	}

	statusCode, body, errs := client.Bytes()
	if statusCode != 200 || errs != nil {
		return nil, fmt.Errorf("unable to get user info")
	}

	googleAuth := &entity.GoogleAuthPayload{}
	if err := json.Unmarshal(body, googleAuth); err != nil {
		return nil, fmt.Errorf("unable to parse user info")
	}
	return googleAuth, nil
}

func GoogleAuthToUser(googleAuth *entity.GoogleAuthPayload) *entity.User {
	return &entity.User{
		FirstName:   googleAuth.GivenName,
		LastName:    googleAuth.FamilyName,
		DisplayName: strings.Split(googleAuth.Email, "@")[0],
		AvatarUrl:   googleAuth.Picture,
		Email:       googleAuth.Email,
		GoogleID:    googleAuth.ID,
	}
}
