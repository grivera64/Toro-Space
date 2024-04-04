package mapper

import (
	"os"
	"strings"

	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/gateway/googleoauth"
)

var (
	adminEmails = map[string]interface{}{
		os.Getenv("ADMIN_EMAIL"): struct{}{},
	}
)

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
