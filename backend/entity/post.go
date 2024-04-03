package entity

import (
	"gorm.io/gorm"
)

type Post struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Content  string `json:"content"`
	Author   *User  `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID uint   `json:"author_id"`
	gorm.Model
}
