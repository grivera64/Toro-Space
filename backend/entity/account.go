package entity

type Account struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	GoogleID  string `json:"google_id"`

	Users []User `json:"users" gorm:"foreignKey:ID"`
}
