package entity

type Post struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`
}
