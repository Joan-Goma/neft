package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	engine "github.com/JoanGTSQ/api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	userPwPPepper = "kedg5b0ays1ekngsg18ruawcekgvcnz6"
	hmacScretKey  = "kedg5b0ays1ekngsg18ruawcekgvcnz6"
	hashCost      = 8
)

type UserDB interface {
	GetAllUsers(pagination Pagination) ([]*User, error)
	Count() (int, error)
}

type UserService interface {
	UserDB
}

var gormUser userGorm
var service userService

func NewUserService(gD *gorm.DB) UserService {
	newUserGorm(gD)
	hmac := engine.NewHMAC(hmacScretKey)
	uv := newUserValidator(&gormUser, hmac)
	service.UserDB = uv
	service.pwResetDB = newPwResetValidator(&pwResetGorm{db: gD}, hmac)

	return &userService{
		UserDB:    uv,
		pwResetDB: newPwResetValidator(&pwResetGorm{db: gD}, hmac),
	}
}

type userService struct {
	UserDB
	pwResetDB pwResetDB
}

func newUserGorm(db *gorm.DB) {
	gormUser.db = db
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (user *User) InitiateReset() (string, error) {
	if err := user.ByID(); err != nil {
		return "", engine.ERR_PSSWD_RESET_TOKEN_DUPLICATED
	}
	pwr := pwReset{
		UserID: user.ID,
	}
	if err := service.pwResetDB.Create(&pwr); err != nil {
		return "", err
	}
	return pwr.TokenHash, nil
}

func (user *User) CompleteReset(token, newPw string) error {
	pwr, err := service.pwResetDB.ByToken(token)
	if err != nil {
		return err
	}
	if time.Since(pwr.CreatedAt) > (2 * time.Hour) {
		return err
	}
	user = &User{
		ID: pwr.UserID,
	}

	err = user.ByID()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPw+userPwPPepper))
	if err == nil {
		return engine.ERR_PSSWD_SAME_RESET
	}
	user.Password = newPw
	err = user.Update()
	if err != nil {
		return err
	}
	service.pwResetDB.Delete(pwr.ID)

	return nil

}

func (user *User) Authenticate() error {
	if user.Email == "" {
		return engine.ERR_MAIL_REQUIRED
	}
	if user.Password == "" {
		return engine.ERR_PSSWD_REQUIRED
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(user.Email) {
		return engine.ERR_MAIL_IS_N0T_VALID
	}
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	err := user.ByEmail()
	if err != nil {
		return engine.ERR_MAIL_NOT_EXIST
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user.Password+userPwPPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return engine.ERR_PSSWD_INCORRECT
		default:
			return err
		}
	}

	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return engine.ERR_NOT_FOUND
	default:
		return err
	}
}

func (ug *userGorm) GetAllUsers(pagination Pagination) ([]*User, error) {

	offset := (pagination.Page - 1) * pagination.Limit
	var users []*User
	err := ug.db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil

}

// BASIC FUNCTIONS
func (user *User) Create() error {
	if err := user.NewUserValidation(); err != nil {
		return err
	}
	return gormUser.db.Create(user).Error
}

func (user *User) Delete() error {
	return gormUser.db.Delete(user).Error
}

func (user *User) Update() error {
	if err := user.ExistentUserValidation(); err != nil {
		return err
	}
	return gormUser.db.Save(user).Error
}

// SEARCH BY
func (user *User) ByID() error {
	db := gormUser.db.Preload("Role").Where("id = ?", user.ID).First(user)
	if err := first(db, user); err != nil {
		return err
	}
	if err := user.CountFollowers(); err != nil {
		return err
	}
	if err := user.CountFollowings(); err != nil {
		return err
	}
	return nil
}

func (user *User) ByEmail() error {
	db := gormUser.db.Preload("Role").Where("email = ?", user.Email)
	if err := first(db, user); err != nil {
		return err
	}
	if err := user.CountFollowers(); err != nil {
		return err
	}
	if err := user.CountFollowings(); err != nil {
		return err
	}
	return nil
}

func (user *User) ByRemember() (*User, error) {
	db := gormUser.db.Preload("Role").Where("remember_hash = ?", user.RememberHash)
	if err := first(db, user); err != nil {
		return nil, err
	}
	if err := user.CountFollowers(); err != nil {
		return nil, err
	}
	if err := user.CountFollowings(); err != nil {
		return nil, err
	}
	return user, nil
}

func (user *User) AssignRole(role int) error {
	return gormUser.db.Model(&User{}).Where("id = ?", user.ID).Update("role_id", role).Error
}
func (user *User) Ban(isBanned bool) error {
	return gormUser.db.Model(&User{}).Where("id = ?", user.ID).Update("banned", isBanned).Error
}

func (ug *userGorm) Count() (int, error) {
	var users int64
	err := ug.db.Table("users").Count(&users).Error
	return int(users), err
}

func (user *User) Follow(friendID uint) error {
	friend := &User{
		ID: friendID,
	}
	err := friend.ByID()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	gormUser.db.Preload("Friends").First(&user, "id = ?", user.ID)
	gormUser.db.Model(&user).Association("Friends").Append(friend)
	return nil
}

func (user *User) Unfollow(friendID uint) error {

	friend := &User{
		ID: friendID,
	}
	err := friend.ByID()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	gormUser.db.Preload("Friends").First(&user, "id = ?", user.ID)
	gormUser.db.Model(&user).Association("Friends").Delete(friend)
	return nil
}

func (user *User) IsFollower(friendID uint) (bool, error) {

	friend := User{}
	friend.ID = friendID
	gormUser.db.Model(&user).Association("Friends").Find(&friend)
	if friend.Email == "" {
		return false, nil
	}
	return true, nil
}
func (user *User) CountFollowings() error {
	user.Following = gormUser.db.Model(&user).Association("Friends").Count()
	return nil
}
func (user *User) CountFollowers() error {
	gormUser.db.Table("friendships").Select("friend_id").Where("friend_id = ?", user.ID).Count(&user.Followers)
	return nil
}

func (user *User) ToString() string {
	return fmt.Sprintf("USER: UserName=%v, FullName=%v, Email=&v", user.UserName, user.FullName, user.Email)
}

type User struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `json:"-" sql:"index"`
	UserName     string     `gorm:"not null" json:"username"`
	FullName     string     `gorm:"not null" json:"full_name"`
	Email        string     `gorm:"not null;unique_index" json:"email"`
	Password     string     `gorm:"-" json:"password"`
	Photo        string     `json:"photo"`
	PasswordHash string     `gorm:"not null" json:"-"`
	Remember     string     `gorm:"-" json:"-"`
	Followers    int        `gorm:"-" json:"followers"`
	Following    int        `gorm:"-" json:"following"`
	Friends      []User     `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id" json:"-"`
	RememberHash string     `gorm:"not null;unique_index" json:"-"`
	RoleID       uint       `gorm:"not null;default: 1" json:"role_id"`
	Role         Role       `gorm:"foreignkey:RoleID" json:"role"`
	Banned       bool       `gorm:"not null; default: false" json:"banned"`
}
