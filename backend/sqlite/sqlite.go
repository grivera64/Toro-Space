package sqlite

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"torospace.csudh.edu/api/entity"
)

type PostParams struct {
	Before   string `json:"before"`
	After    string `json:"after"`
	PageSize int    `json:"page_size"`
}

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

func (db *DB) GetUserByID(id uint) (*entity.User, error) {
	db.Lock()
	defer db.Unlock()

	user := &entity.User{}
	err := db.gormDB.First(user, "id = ?", id).Error
	return user, err
}

func (db *DB) AddPost(post *entity.Post) error {
	db.Lock()
	defer db.Unlock()

	return db.gormDB.Create(post).Error
}

func (db *DB) GetPosts(params *PostParams) ([]*entity.Post, error) {
	db.Lock()
	defer db.Unlock()

	// By default, provide latest 10 posts
	if params == nil {
		params = &PostParams{
			PageSize: 10,
		}
	}

	var posts []*entity.Post
	query := db.gormDB.Preload("Author").
		Order("created_at DESC")

	if params.Before != "" {
		query = query.Where("id < ?", params.Before)
	}

	if params.After != "" {
		query = query.Where("id > ?", params.After)
	}

	err := query.Limit(params.PageSize).
		Find(&posts).
		Error
	return posts, err
}

func (db *DB) GetPostsByUserID(id uint, params *PostParams) ([]*entity.Post, error) {
	db.Lock()
	defer db.Unlock()

	// By default, provide latest 10 posts
	if params == nil {
		params = &PostParams{
			PageSize: 10,
		}
	}

	var posts []*entity.Post
	query := db.gormDB.Preload("Users").
		Order("created_at DESC").
		Where("author_id = ?", id)

	if params.Before != "" {
		query = query.Where("created_at < ?", params.Before)
	}

	if params.After != "" {
		query = query.Where("created_at > ?", params.After)
	}

	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	err := query.Limit(params.PageSize).
		Find(&posts).
		Error
	return posts, err
}

func (db *DB) CreateTopic(topic *entity.Topic) error {
	db.Lock()
	defer db.Unlock()

	return db.gormDB.Create(topic).Error
}

func (db *DB) GetTopicByName(name string) (*entity.Topic, error) {
	db.Lock()
	defer db.Unlock()

	topic := &entity.Topic{}
	err := db.gormDB.First(&topic, "name = ?", name).Error
	return topic, err
}
