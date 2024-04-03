package googleoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

// GoogleAuthGateway is a gateway for Google authentication.
type GoogleOauthGateway interface {
	GetAuthUrl() string
	GetToken(c context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(token string) (*GoogleUser, error)
}

type googleOauthGatewayV2 struct {
	baseUrl     string
	oauthConfig *oauth2.Config
}

func NewV2() GoogleOauthGateway {
	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("G_CLIENT_ID"),
		ClientSecret: os.Getenv("G_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("G_REDIRECT"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	return &googleOauthGatewayV2{
		baseUrl:     "https://www.googleapis.com/oauth2/v1",
		oauthConfig: oauthConfig,
	}
}

func (g *googleOauthGatewayV2) GetAuthUrl() string {
	// Todo add state and opts
	return g.oauthConfig.AuthCodeURL("state")
}

func (g *googleOauthGatewayV2) GetToken(c context.Context, authCode string) (*oauth2.Token, error) {
	return g.oauthConfig.Exchange(c, authCode)
}

func (g *googleOauthGatewayV2) GetUserInfo(token string) (*GoogleUser, error) {
	client := fiber.Get(fmt.Sprintf("%s/userinfo?alt=json&access_token=%s", g.baseUrl, token))
	if client == nil {
		return nil, fmt.Errorf("unable to get user info")
	}

	statusCode, body, errs := client.Bytes()
	if statusCode != 200 || errs != nil {
		return nil, fmt.Errorf("unable to get user info")
	}

	googleAuth := &GoogleUser{}
	if err := json.Unmarshal(body, googleAuth); err != nil {
		return nil, fmt.Errorf("unable to parse user info")
	}
	return googleAuth, nil
}
