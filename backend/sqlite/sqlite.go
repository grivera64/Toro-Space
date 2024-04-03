package sqlite

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"torospace.csudh.edu/api/entity"
)

type DB struct {
	gormDB *gorm.DB
	sync.Mutex
}

func NewDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("torospace.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.Account{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Post{})
	return &DB{gormDB: db}, nil
}

func (db *DB) AddAccount(account *entity.Account) error {
	db.Lock()
	defer db.Unlock()

	return db.gormDB.Create(account).Error
}

func (db *DB) AddAccountUser(account *entity.Account, user *entity.User) error {
	db.Lock()
	defer db.Unlock()

	account.Users = append(account.Users, *user)
	return db.gormDB.Save(account).Error
}

func (db *DB) GetAccountByID(id uint) (*entity.Account, error) {
	db.Lock()
	defer db.Unlock()

	user := &entity.Account{}
	err := db.gormDB.Preload("Users").First(user, "id = ?", id).Error
	return user, err
}

func (db *DB) GetAccountByGoogleID(id string) (*entity.Account, error) {
	db.Lock()
	defer db.Unlock()

	user := &entity.Account{}
	err := db.gormDB.First(user, "google_id = ?", id).Error
	return user, err
}
