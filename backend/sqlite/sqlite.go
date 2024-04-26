package sqlite

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"torospace.csudh.edu/api/entity"
)

type PostParams struct {
	Before      string `json:"before"`
	After       string `json:"after"`
	PageSize    int    `json:"page_size"`
	SearchQuery string `json:"search_query"`
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
	db.AutoMigrate(&entity.Topic{})
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

	account := &entity.Account{}
	err := db.gormDB.Preload("Users").First(account, "id = ?", id).Error
	return account, err
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

	return db.gormDB.Preload("Topics").Create(post).Error
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
	query := db.gormDB.Preload("LikedBy").Preload("Author").Preload("Topics").
		Order("created_at DESC")

	if params.Before != "" {
		query = query.Where("id < ?", params.Before)
	}

	if params.After != "" {
		query = query.Where("id > ?", params.After)
	}

	if params.SearchQuery != "" {
		query = query.Where("content LIKE ?", fmt.Sprintf("%%%s%%", params.SearchQuery))
	}

	err := query.Limit(params.PageSize).
		Find(&posts).
		Error
	return posts, err
}

func (db *DB) GetPost(postID uint) (*entity.Post, error) {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.Preload("LikedBy").Preload("Author").Preload("Topics").First(post, "id = ?", postID).Error
	return post, err
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
	query := db.gormDB.Preload("LikedBy").Preload("Author").
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

func (db *DB) AddLikeToPost(postID uint, user *entity.User) error {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.Preload("LikedBy").Preload("Author").First(post, "id = ?", postID).Error
	if err != nil {
		return err
	}

	if post.LikedBy == nil {
		post.LikedBy = []entity.User{}
	}
	post.Likes = len(post.LikedBy)

	for _, u := range post.LikedBy {
		if u.ID == user.ID {
			log.Println("User already liked post")
			return nil
		}
	}

	post.LikedBy = append(post.LikedBy, *user)
	post.Likes = len(post.LikedBy)
	return db.gormDB.Preload("LikedBy").Preload("Author").Save(post).Error
}

func (db *DB) RemoveLikeFromPost(postID uint, user *entity.User) error {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.Preload("LikedBy").Preload("Author").First(post, "id = ?", postID).Error
	if err != nil {
		return err
	}
	if post.LikedBy == nil {
		return nil
	}
	if post.Likes == 0 {
		return nil
	}

	indexMatch := -1
	for index, u := range post.LikedBy {
		if u.ID == user.ID {
			indexMatch = index
			break
		}
	}

	if indexMatch != -1 {
		toRemove := post.LikedBy[indexMatch]
		log.Println(post.LikedBy)
		if err := db.gormDB.Model(post).Association("LikedBy").Delete(toRemove); err != nil {
			log.Println("Error deleting like: ", err)
			return err
		}
		post.Likes = len(post.LikedBy)
		err := db.gormDB.Preload("LikedBy").Preload("Author").Save(post).Error
		if err != nil {
			return err
		}
		log.Println(post.LikedBy)
		return nil
	}

	return fmt.Errorf("user %d has not liked post %d", user.ID, postID)
}

func (db *DB) GetPostLikesByID(postID uint) ([]entity.User, error) {
	db.Lock()
	defer db.Unlock()

	var post entity.Post
	err := db.gormDB.Preload("LikedBy").Preload("Author").First(post).Error
	return post.LikedBy, err
}

func (db *DB) CreateTopic(topic *entity.Topic) error {
	db.Lock()
	defer db.Unlock()

	return db.gormDB.Create(topic).Error
}

func (db *DB) GetTopics() ([]*entity.Topic, error) {
	db.Lock()
	defer db.Unlock()

	var topics []*entity.Topic
	err := db.gormDB.Find(&topics).Error
	return topics, err
}

func (db *DB) GetTopicByName(name string) (*entity.Topic, error) {
	db.Lock()
	defer db.Unlock()

	topic := &entity.Topic{}
	err := db.gormDB.First(topic, "name = ?", name).Error
	return topic, err
}
