package entity

import (
	"encoding/gob"

	"torospace.csudh.edu/api/util"
)

type Role string

var (
	RoleOrganization Role = "organization"
	RoleStudent      Role = "student"
	RoleAdmin        Role = "admin"
)

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	DisplayName string `json:"display_name"`
	AvatarUrl   string `json:"avatar_url"`
	Role        Role   `json:"role"`
}

func init() {
	gob.Register(Role(""))
}

func (u User) LessThan(other util.Comparable) bool {
	otherUser, ok := other.(User)
	if !ok {
		return false
	}
	return u.ID < otherUser.ID
}

func (u User) GreaterThan(other util.Comparable) bool {
	otherUser, ok := other.(User)
	if !ok {
		return false
	}
	return u.ID > otherUser.ID
}

func (u User) EqualTo(other util.Comparable) bool {
	otherUser, ok := other.(User)
	if !ok {
		return false
	}
	return u.ID == otherUser.ID
}
