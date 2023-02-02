package models

import (
	engine "github.com/JoanGTSQ/api"
	_ "github.com/lib/pq"
)

func (post *Post) Create() error {
	err := DBCONNECTION.Create(post).Error
	if err != nil {
		return err
	}
	return nil
}

func (post *Post) Delete() error {
	return DBCONNECTION.Delete(&post).Error
}

func (post *Post) Update() error {
	return DBCONNECTION.Save(post).Error
}

// SEARCH BY ID

func (post *Post) ByID() error {
	if err := DBCONNECTION.Where("id = ?", post.ID).Preload("User").First(post).Error; err != nil {
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

// AllPosts search in the entire database with the pagination and save all the results
func (mp *MultiplePost) AllPosts(userID uint) error {
	offset := (mp.Pagination.Page - 1) * mp.Pagination.Limit
	err := DBCONNECTION.Offset(offset).Limit(mp.Pagination.Limit).Order(mp.Pagination.Sort).Where("user_id = ?", userID).Preload("User").Find(&mp.Posts).Error
	for _, p := range mp.Posts {
		if err := p.User.CountFollowers(); err != nil {
			return err
		}
		if err := p.CountLikes(); err != nil {
			return err
		}

		if err := p.GetComments(); err != nil {
			return err
		}
	}
	return err
}

// Count count all the matches posible in the database
func (mp *MultiplePost) Count() error {
	err := DBCONNECTION.Table("posts").Count(&mp.Quantity).Error
	return err
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
	post.LikesReceived = DBCONNECTION.Model(&post).Association("Likes").Count()
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
func (post *Post) Unlike(friendID uint) error {
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

	DBCONNECTION.Create(comment)

	DBCONNECTION.First(&post, "id = ?", post.ID)
	return DBCONNECTION.Model(&post).Association("post_comments").Append(comment).Error
}

// Uncomment delete the comment and Delete from the association
func (post *Post) Uncomment(comment *Comment) error {
	DBCONNECTION.First(comment, "id = ?", comment.ID)
	DBCONNECTION.Delete(comment)
	DBCONNECTION.Preload("Friends").First(&post, "id = ?", post.ID)
	DBCONNECTION.Model(&post).Association("post_comments").Delete(comment)
	return nil
}

// GetComments Get all comments from a post
func (post *Post) GetComments() error {

	DBCONNECTION.
		Preload("PostComments").
		Preload("PostComments.User").
		First(&post, "id = ?", post.ID)
	return nil
}

type MultiplePost struct {
	Pagination Pagination
	Posts      []*Post
	Quantity   int64
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
