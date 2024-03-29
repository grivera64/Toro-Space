package entity

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	AvatarUrl   string `json:"avatar_url"`
	Email       string `json:"email"`
	GoogleID    string `json:"google_id"`
}
