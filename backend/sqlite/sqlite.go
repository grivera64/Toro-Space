package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"torospace.csudh.edu/api/entity"
)

type DB struct {
	gormDB *gorm.DB
}

func NewDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("torospace.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.User{})
	return &DB{gormDB: db}, nil
}

func (db *DB) AddUser(user *entity.User) error {
	return db.gormDB.Create(user).Error
}

func (db *DB) GetUserByID(id uint) (*entity.User, error) {
	user := &entity.User{}
	err := db.gormDB.First(user, "id = ?", id).Error
	return user, err
}

func (db *DB) GetUserByGoogleID(id string) (*entity.User, error) {
	user := &entity.User{}
	err := db.gormDB.First(user, "google_id = ?", id).Error
	return user, err
}
