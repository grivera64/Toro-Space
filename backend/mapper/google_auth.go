package mapper

import (
	"strings"

	"torospace.csudh.edu/api/entity"
	"torospace.csudh.edu/api/gateway/googleoauth"
)

func GoogleUserToAccount(googleUser *googleoauth.GoogleUser) *entity.Account {
	return &entity.Account{
		FirstName: googleUser.GivenName,
		LastName:  googleUser.FamilyName,
		Email:     googleUser.Email,
		GoogleID:  googleUser.ID,
		Users: []entity.User{
			{
				DisplayName: strings.Split(googleUser.Email, "@")[0],
				AvatarUrl:   googleUser.Picture,
				Role:        entity.RoleStudent,
			},
		},
	}
}
