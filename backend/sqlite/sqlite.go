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
	GetHidden   bool   `json:"get_hidden"`
}

type PostsResult struct {
	Posts     []*entity.Post `json:"posts"`
	Count     int            `json:"count"`
	HasBefore bool           `json:"has_before"`
	HasAfter  bool           `json:"has_after"`
}

type OrganizationParams struct {
	Before      string `json:"before"`
	After       string `json:"after"`
	PageSize    int    `json:"page_size"`
	SearchQuery string `json:"search_query"`
}

type OrganizationsResult struct {
	Organizations []*entity.User `json:"organizations"`
	Count         int            `json:"count"`
	HasBefore     bool           `json:"has_before"`
	HasAfter      bool           `json:"has_after"`
}

type TopicParams struct {
	Before      string `json:"before"`
	After       string `json:"after"`
	PageSize    int    `json:"page_size"`
	SearchQuery string `json:"search_query"`
}

type TopicsResult struct {
	Topics    []*entity.Topic `json:"topics"`
	Count     int             `json:"count"`
	HasBefore bool            `json:"has_before"`
	HasAfter  bool            `json:"has_after"`
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

func (db *DB) GetPosts(params *PostParams) (*PostsResult, error) {
	db.Lock()
	defer db.Unlock()

	// By default, provide latest 10 posts
	if params == nil {
		params = &PostParams{
			PageSize:  10,
			GetHidden: false,
		}
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	query := db.gormDB.Model(&entity.Post{}).Preload("LikedBy").Preload("Author").Preload("Topics").
		Order("created_at DESC")

	if !params.GetHidden {
		query = query.Where("hidden <> ?", true)
	}

	if params.SearchQuery != "" {
		var matchingByTopic []uint64
		var matchingByAuthor []uint64
		searchQuery := fmt.Sprintf("%%%s%%", params.SearchQuery)
		db.gormDB.
			Table("posts").
			Select("posts.id").
			Joins("JOIN post_topics ON posts.id = post_topics.post_id").
			Joins("JOIN topics ON post_topics.topic_id = topics.id").
			Where("topics.name LIKE ?", searchQuery).
			Distinct().
			Pluck("posts.id", &matchingByTopic)

		db.gormDB.
			Table("posts").
			Select("posts.id").
			Joins("JOIN users ON posts.author_id = users.id").
			Where("users.display_name LIKE ?", searchQuery).
			Distinct().
			Pluck("posts.id", &matchingByAuthor)

		query = query.Where("(content LIKE (?)) OR (posts.id IN (?)) OR (posts.id IN (?))", fmt.Sprintf("%%%s%%", params.SearchQuery), matchingByTopic, matchingByAuthor)
	}

	var totalCount int64
	var newestPost entity.Post
	var oldestPost entity.Post
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	if err := query.First(&newestPost).Error; err != nil {
		return nil, err
	}
	if err := query.Offset(int(totalCount - 1)).Limit(1).Find(&oldestPost).Error; err != nil {
		return nil, err
	}

	query = query.Offset(-1).Limit(-1)

	if params.Before != "" {
		query = query.Where("id < ?", params.Before)
	}

	if params.After != "" {
		query = query.Where("id > ?", params.After)
	}

	var posts []*entity.Post
	err := query.Limit(params.PageSize).
		Find(&posts).
		Error

	result := &PostsResult{}
	if len(posts) > 0 {
		result.HasBefore = posts[len(posts)-1].ID > oldestPost.ID
		result.HasAfter = posts[0].ID < newestPost.ID
	} else {
		result.HasBefore = false
		result.HasAfter = false
	}

	result.Posts = posts
	result.Count = len(posts)

	return result, err
}

func (db *DB) GetPost(postID uint) (*entity.Post, error) {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.Preload("LikedBy").Preload("Author").Preload("Topics").First(post, "id = ?", postID).Error
	return post, err
}

func (db *DB) DeletePost(postID uint) error {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.First(post, "id = ?", postID).Error
	if err != nil {
		return err
	}
	return db.gormDB.Delete(post).Error
}

func (db *DB) HidePost(postID uint) error {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.First(post, "id = ?", postID).Error
	if err != nil {
		return err
	}

	if post.Hidden {
		return nil
	}
	post.Hidden = true
	return db.gormDB.Save(post).Error
}

func (db *DB) UnhidePost(postID uint) error {
	db.Lock()
	defer db.Unlock()

	post := &entity.Post{}
	err := db.gormDB.First(post, "id = ?", postID).Error
	if err != nil {
		return err
	}

	if !post.Hidden {
		return nil
	}
	post.Hidden = false
	return db.gormDB.Save(post).Error
}

func (db *DB) GetPostsByOrganization(id uint, params *PostParams) (*PostsResult, error) {
	db.Lock()
	defer db.Unlock()

	// By default, provide latest 10 posts
	if params == nil {
		params = &PostParams{
			PageSize: 10,
		}
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	{
		user := &entity.User{}
		db.gormDB.First(user, "id = ?", id)
		if user.Role != entity.RoleOrganization {
			return nil, fmt.Errorf("user's role is not organization")
		}
	}

	query := db.gormDB.Preload("LikedBy").Preload("Author").Preload("Topics").
		Order("created_at DESC").
		Where("author_id = ?", id)

	if !params.GetHidden {
		query = query.Where("hidden = ?", false)
	}

	if params.SearchQuery != "" {
		var matchingIDs []uint64
		searchQuery := fmt.Sprintf("%%%s%%", params.SearchQuery)
		db.gormDB.
			Table("posts").
			Select("posts.id").
			Joins("JOIN post_topics ON posts.id = post_topics.post_id").
			Joins("JOIN topics ON post_topics.topic_id = topics.id").
			Where("topics.name LIKE ?", searchQuery).
			Distinct().
			Pluck("posts.id", &matchingIDs)

		query = query.Where("(content LIKE (?)) OR (posts.id IN (?))", fmt.Sprintf("%%%s%%", params.SearchQuery), matchingIDs)
	}

	var totalCount int64
	var newestPost entity.Post
	var oldestPost entity.Post
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	if err := query.First(&newestPost).Error; err != nil {
		return nil, err
	}
	if err := query.Offset(int(totalCount - 1)).Limit(1).Find(&oldestPost).Error; err != nil {
		return nil, err
	}

	query = query.Offset(-1).Limit(-1)

	if params.Before != "" {
		query = query.Where("id < ?", params.Before)
	}

	if params.After != "" {
		query = query.Where("id > ?", params.After)
	}

	result := &PostsResult{}
	var posts []*entity.Post
	err := query.Limit(params.PageSize).
		Find(&posts).
		Error

	if len(posts) > 0 {
		result.HasBefore = posts[len(posts)-1].ID > oldestPost.ID
		result.HasAfter = posts[0].ID < newestPost.ID
	} else {
		result.HasBefore = false
		result.HasAfter = false
	}

	result.Posts = posts
	result.Count = len(posts)

	return result, err
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

func (db *DB) GetTopics(params *TopicParams) (*TopicsResult, error) {
	db.Lock()
	defer db.Unlock()

	if params == nil {
		params = &TopicParams{
			PageSize: 10,
		}
	}

	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	var topics []*entity.Topic
	query := db.gormDB.Model(&entity.Topic{})

	if params.SearchQuery != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", params.SearchQuery))
	}

	var totalCount int64
	var firstTopic entity.Topic
	var lastTopic entity.Topic
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	if err := query.First(&firstTopic).Error; err != nil {
		return nil, err
	}
	if err := query.Offset(int(totalCount - 1)).Limit(1).Find(&lastTopic).Error; err != nil {
		return nil, err
	}

	query = query.Offset(-1).Limit(-1)

	if params.Before != "" {
		query = query.Where("id < ?", params.Before)
	}
	if params.After != "" {
		query = query.Where("id > ?", params.After)
	}

	result := &TopicsResult{}
	if len(topics) > 0 {
		result.HasBefore = topics[len(topics)-1].ID > lastTopic.ID
		result.HasAfter = topics[0].ID < firstTopic.ID
	} else {
		result.HasBefore = false
		result.HasAfter = false
	}

	err := query.Limit(params.PageSize).
		Find(&topics).Error

	result.Topics = topics
	result.Count = len(topics)

	return result, err
}

func (db *DB) GetTopicByName(name string) (*entity.Topic, error) {
	db.Lock()
	defer db.Unlock()

	topic := &entity.Topic{}
	err := db.gormDB.First(topic, "name = ?", name).Error
	return topic, err
}

func (db *DB) GetOrganizations(params *OrganizationParams) (*OrganizationsResult, error) {
	db.Lock()
	defer db.Unlock()

	if params == nil {
		params = &OrganizationParams{
			PageSize: 10,
		}
	}

	var organizations []*entity.User
	query := db.gormDB.Model(&entity.User{}).
		Where("role = ?", "organization")

	if params.SearchQuery != "" {
		query = query.Where("display_name LIKE ?", fmt.Sprintf("%%%s%%", params.SearchQuery))
	}

	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	var totalCount int64
	var firstOrganization entity.User
	var lastOrganization entity.User
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	if err := query.First(&firstOrganization).Error; err != nil {
		return nil, err
	}
	if err := query.Offset(int(totalCount - 1)).Limit(1).Find(&lastOrganization).Error; err != nil {
		return nil, err
	}

	query = query.Offset(-1).Limit(-1)

	if params.Before != "" {
		query = query.Where("id < ?", params.Before)
	}
	if params.After != "" {
		query = query.Where("id > ?", params.After)
	}

	result := &OrganizationsResult{}
	if len(organizations) > 0 {
		result.HasBefore = organizations[len(organizations)-1].ID > lastOrganization.ID
		result.HasAfter = organizations[0].ID < firstOrganization.ID
	} else {
		result.HasBefore = false
		result.HasAfter = false
	}

	err := query.Limit(params.PageSize).
		Find(&organizations).Error

	result.Organizations = organizations
	result.Count = len(organizations)

	return result, err
}

func (db *DB) GetOrganization(id uint) (*entity.User, error) {
	db.Lock()
	defer db.Unlock()

	organization := &entity.User{}
	err := db.gormDB.First(organization, "id = ?", id).Error
	return organization, err
}
