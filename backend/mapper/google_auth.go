package mapper

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/gateway/googleoauth"
)

var (
	adminEmails = map[string]interface{}{}
)

func init() {
	godotenv.Load()
	adminEmail, ok := os.LookupEnv("ADMIN_EMAIL")
	if !ok {
		panic("ADMIN_EMAIL environment variable is not set")
	}
	adminEmails[adminEmail] = struct{}{}
}

func GoogleUserToAccount(googleUser *googleoauth.GoogleUser) *entity.Account {
	role := entity.RoleStudent
	if _, ok := adminEmails[googleUser.Email]; ok {
		role = entity.RoleAdmin
	}

	return &entity.Account{
		FirstName: googleUser.GivenName,
		LastName:  googleUser.FamilyName,
		Email:     googleUser.Email,
		GoogleID:  googleUser.ID,
		Users: []entity.User{
			{
				DisplayName: strings.Split(googleUser.Email, "@")[0],
				AvatarUrl:   googleUser.Picture,
				Role:        role,
			},
		},
	}
}
