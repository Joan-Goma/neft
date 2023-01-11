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

func NewServices(connectionInfo string, logMode bool) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(logMode)

	return &Services{
		User:     NewUserService(db),
		Spot:     NewSpotService(db),
		Category: NewCategoryService(db),
		Post:     NewPostService(db),
		Device:   NewDeviceService(db),
		Role:     NewRoleService(db),
		Message:  NewMessageService(db),
		db:       db,
	}, nil
}

type Services struct {
	Spot     SpotService
	User     UserService
	Category CategoryService
	Post     PostService
	Device   DeviceService
	Message  MessageService
	Role     RoleService
	db       *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&UserMessage{}, &Category{}, &pwReset{}, &Role{}, &User{}, &Spot{}, &Comment{}, &Post{},
		&Device{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) DestructiveStatic() error {
	if err := s.db.DropTableIfExists().Error; err != nil {
		return err
	}
	if err := s.AutoMigrate(); err != nil {
		return err
	}
	r1 := &Role{
		RoleName: "User",
	}
	r2 := &Role{
		RoleName:    "Moderator",
		CanBanUsers: true,
	}
	r3 := &Role{
		RoleName:    "Moderator",
		CanBanUsers: true,
	}
	s.db.Model(&Role{}).Create(r1)
	s.db.Model(&Role{}).Create(r2)
	s.db.Model(&Role{}).Create(r3)
	return nil
}

func (s *Services) AutoMigrate() error {

	if err := s.db.AutoMigrate(&User{}, &UserMessage{}, &Role{}, &pwReset{}, &Category{}, &Spot{}, &Comment{}, &Post{},
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
