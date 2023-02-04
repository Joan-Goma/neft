package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func NewServices(connectionInfo string, logMode bool) error {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return err
	}

	db.LogMode(logMode)
	DBCONNECTION = db
	NewUserService(db)
	NewCategoryService(db)
	NewDeviceService(db)
	return nil
}

var DBCONNECTION *gorm.DB

type Services struct {
	User     UserService
	Category CategoryService
	Device   DeviceService
	db       *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func DestructiveReset() error {
	if err := DBCONNECTION.DropTableIfExists(&UserMessage{}, &Category{}, &pwReset{}, &Role{}, &User{}, &Spot{}, &Comment{},
		&Device{}).Error; err != nil {
		return err
	}
	return AutoMigrate()
}

func AutoMigrate() error {

	if err := DBCONNECTION.AutoMigrate(&User{}, &UserMessage{}, &Role{}, &pwReset{}, &Category{}, &Comment{}, &Spot{},
		&Device{}).Error; err != nil {
		return err
	}
	return nil
}

type NeftModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
