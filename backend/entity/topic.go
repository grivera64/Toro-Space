package entity

type Topic struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
