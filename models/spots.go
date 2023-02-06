package models

import (
	engine "github.com/JoanGTSQ/api"
	_ "github.com/lib/pq"
)

func (spot *Spot) Create() error {
	if err := spot.ValidateNewSpot(); err != nil {
		return err
	}
	return DBCONNECTION.Create(spot).Error
}

func (spot *Spot) Delete() error {
	return DBCONNECTION.Delete(&spot).Error
}

func (spot *Spot) Update() error {
	if err := spot.ValidateNewSpot(); err != nil {
		return err
	}
	return DBCONNECTION.Save(spot).Error
}

// SEARCH BY ID

func (spot *Spot) ByID() error {
	if err := DBCONNECTION.
		Where("id = ?", spot.ID).
		First(&spot).
		Error; err != nil {
		return err
	}

	if err := spot.CountLikes(); err != nil {
		return err
	}

	if err := spot.GetComments(); err != nil {
		return err
	}
	return nil
}

func (ms *MultipleSpots) AllSpots() error {
	offset := (ms.Pagination.Page - 1) * ms.Pagination.Limit
	err := DBCONNECTION.
		Offset(offset).
		Limit(ms.Pagination.Limit).
		Order(ms.Pagination.Sort).
		Find(&ms.Spots).
		Error
	for _, s := range ms.Spots {
		if err := s.CountLikes(); err != nil {
			return err
		}

		if err := s.GetComments(); err != nil {
			return err
		}
	}
	return err
}

func (ms *MultipleSpots) Count() error {
	var i int64
	err := DBCONNECTION.Table("spots").Count(&i).Error
	ms.Quantity = i
	return err
}

func (spot *Spot) CountLikes() error {
	spot.LikesReceived = DBCONNECTION.Model(&spot).Association("Likes").Count()
	return nil
}

func (spot *Spot) Like(likerID uint) error {
	liker := &User{
		ID: likerID,
	}
	err := liker.ByID()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	gormUser.db.Preload("Likes").First(&spot, "id = ?", spot.ID)
	gormUser.db.Model(&spot).Association("Likes").Append(liker)
	return nil
}

// Unlike delete the comment and Delete from the association
func (spot *Spot) Unlike(friendID uint) error {
	liker := &User{
		ID: friendID,
	}
	err := liker.ByID()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	gormUser.db.Preload("Likes").First(&spot, "id = ?", spot.ID)
	gormUser.db.Model(&spot).Association("Likes").Delete(liker)
	return nil
}

// Comment create the comment and add append it with the spot ID
func (spot *Spot) Comment(comment *Comment) error {

	DBCONNECTION.Create(comment)

	DBCONNECTION.First(&spot, "id = ?", spot.ID)
	return DBCONNECTION.Model(&spot).Association("spot_comments").Append(comment).Error
}

// Uncomment delete the comment and Delete from the association
func (spot *Spot) Uncomment(comment *Comment) error {
	DBCONNECTION.First(comment, "id = ?", comment.ID)
	DBCONNECTION.Delete(comment)
	DBCONNECTION.Preload("Friends").First(&spot, "id = ?", spot.ID)
	DBCONNECTION.Model(&spot).Association("spot_comments").Delete(comment)
	return nil
}

// GetComments Get all comments from a spot
func (spot *Spot) GetComments() error {

	DBCONNECTION.
		Preload("SpotComments").
		Preload("SpotComments.User").
		First(&spot, "id = ?", spot.ID)
	return nil
}

type MultipleSpots struct {
	Pagination Pagination
	Spots      []*Spot
	Quantity   int64
}

type Spot struct {
	NeftModel
	Description   string    `gorm:"not null" json:"description"`
	Coordinates   string    `gorm:"not null" json:"coordinates,omitempty"`
	Visible       bool      `gorm:"default: true" json:"visible"`
	UserID        uint      `gorm:"not null" json:"-"`
	User          User      `gorm:"foreignkey:UserID; preload: true" json:"creator,omitempty"`
	LikesReceived int       `gorm:"-" json:"likes,omitempty"`
	Likes         []User    `gorm:"many2many:likes;association_jointable_foreignkey:user_id" json:"-"`
	SpotComments  []Comment `gorm:"many2many:spot_comments;association_jointable_foreignkey:comment_id; preload: true" json:"comments,omitempty"`
}

type Comment struct {
	NeftModel
	Text          string `gorm:"not null" json:"text"`
	CommentatorID uint   `gorm:"not null; preload: true" json:"-"`
	User          User   `gorm:"foreignkey:CommentatorID; preload: true" json:"commentator"`
}
