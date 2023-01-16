package models

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type PostDB interface {
	AllPosts(pagination Pagination, userID uint) ([]*Post, error)
	Count() (int, error)
}

type PostService interface {
	PostDB
}

var gormPost postGorm

func newPostGorm(db *gorm.DB) (*postGorm, error) {
	gormPost.db = db
	return &postGorm{
		db: db,
	}, nil
}
func NewPostService(gD *gorm.DB) PostService {
	ug, err := newPostGorm(gD)
	if err != nil {
		return nil
	}
	return &postService{
		PostDB: ug,
	}
}

type postService struct {
	PostDB
}

var _ PostDB = &postGorm{}

type postGorm struct {
	db *gorm.DB
}

func (post *Post) Create() error {
	err := gormPost.db.Create(post).Error
	if err != nil {
		return err
	}
	return nil
}

func (post *Post) Delete() error {
	return gormPost.db.Delete(&post).Error
}

func (post *Post) Update() error {
	return gormPost.db.Save(post).Error
}

// SEARCH BY ID

func (post *Post) ByID() error {
	if err := gormPost.db.Where("id = ?", post.ID).Preload("User").First(post).Error; err != nil {
		return engine.ERR_NOT_FOUND
	}

	if err := post.User.CountFollowings(); err != nil {
		return err
	}

	if err := post.User.CountFollowers(); err != nil {
		return err
	}
	if err := post.CountLikes(); err != nil {
		return err
	}

	if err := post.GetComments(); err != nil {
		return err
	}
	return nil
}

func (ug *postGorm) AllPosts(pagination Pagination, userID uint) ([]*Post, error) {
	var post []*Post
	offset := (pagination.Page - 1) * pagination.Limit
	err := ug.db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort).Where("user_id = ?", userID).Preload("User").Find(&post).Error
	for _, p := range post {
		if err := p.User.CountFollowers(); err != nil {
			return nil, err
		}
		if err := p.CountLikes(); err != nil {
			return nil, err
		}

		if err := p.GetComments(); err != nil {
			return nil, err
		}
	}
	return post, err
}

func (tg *postGorm) Count() (int, error) {
	var posts int64
	err := tg.db.Table("posts").Count(&posts).Error
	return int(posts), err
}

// Ban a post
func (post *Post) Ban() error {
	engine.Debug.Println("New post banned ID:", post.ID)
	return gormUser.db.Model(&User{}).Where("id = ?", post.ID).Update("visible", false).Error
}

// Unban a post
func (post *Post) Unban() error {
	engine.Debug.Println("New post unbanned ID:", post.ID)
	return gormUser.db.Model(&User{}).Where("id = ?", post.ID).Update("visible", true).Error
}

func (post *Post) CountLikes() error {
	post.LikesReceived = gormPost.db.Model(&post).Association("Likes").Count()
	return nil
}

func (post *Post) Like(friendID uint) error {
	liker := &User{
		ID: friendID,
	}
	err := liker.ByID()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	gormUser.db.Preload("Likes").First(&post, "id = ?", post.ID)
	gormUser.db.Model(&post).Association("Likes").Append(liker)
	return nil
}

// Unlike delete the comment and Delete from the association
func (post Post) Unlike(friendID uint) error {
	liker := &User{
		ID: friendID,
	}
	err := liker.ByID()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	gormUser.db.Preload("Likes").First(&post, "id = ?", post.ID)
	gormUser.db.Model(&post).Association("Likes").Delete(liker)
	return nil
}

// Comment create the comment and add append it with the post ID
func (post *Post) Comment(comment *Comment) error {

	gormPost.db.Create(comment)

	gormPost.db.First(&post, "id = ?", post.ID)
	return gormPost.db.Model(&post).Association("post_comments").Append(comment).Error
}

// Uncomment delete the comment and Delete from the association
func (post Post) Uncomment(comment *Comment) error {
	gormPost.db.First(comment, "id = ?", comment.ID)
	gormPost.db.Delete(comment)
	gormPost.db.Preload("Friends").First(&post, "id = ?", post.ID)
	gormPost.db.Model(&post).Association("post_comments").Delete(comment)
	return nil
}

// GetComments Get all comments from a post
func (post *Post) GetComments() error {

	gormPost.db.
		Preload("PostComments").
		Preload("PostComments.User").
		First(&post, "id = ?", post.ID)
	return nil
}

type Post struct {
	NeftModel
	Description   string    `json:"description"`
	Visible       bool      `gorm:"default: true" json:"visible"`
	Photo         string    `gorm:"not null" json:"photo"`
	UserID        uint      `gorm:"not null" json:"-"`
	User          User      `gorm:"foreignkey:UserID" json:"creator"`
	LikesReceived int       `gorm:"-" json:"likes"`
	Likes         []User    `gorm:"many2many:likes;association_jointable_foreignkey:user_id" json:"-"`
	PostComments  []Comment `gorm:"many2many:post_comments;association_jointable_foreignkey:comment_id; preload: true" json:"comments"`
}

type Comment struct {
	NeftModel
	Text          string `gorm:"not null" json:"text"`
	CommentatorID uint   `gorm:"not null; preload: true" json:"-"`
	User          User   `gorm:"foreignkey:CommentatorID; preload: true" json:"commentator"`
}
