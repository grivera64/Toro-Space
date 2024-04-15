package entity

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Content  string  `json:"content"`
	Author   User    `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID uint    `json:"author_id"`
	Topics   []Topic `json:"topics" gorm:"many2many:post_topics;"`

	LikedBy []*User `json:"liked_by" gorm:"many2many:post_users"`
	Likes   int     `json:"likes"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
