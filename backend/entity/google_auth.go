package entity

type GoogleAuthPayload struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func (payload *GoogleAuthPayload) GetEmail() string {
	return payload.Email
}

func (payload *GoogleAuthPayload) GetName() string {
	return payload.Name
}

func (payload *GoogleAuthPayload) GetFirstName() string {
	return payload.GivenName
}

func (payload *GoogleAuthPayload) GetLastName() string {
	return payload.FamilyName
}

func (payload *GoogleAuthPayload) GetAvatarUrl() string {
	return payload.Picture
}

func (payload *GoogleAuthPayload) GetLocale() string {
	return payload.Locale
}

func (payload *GoogleAuthPayload) IsEmailVerified() bool {
	return payload.EmailVerified
}

func (payload *GoogleAuthPayload) GetID() string {
	return payload.ID
}
