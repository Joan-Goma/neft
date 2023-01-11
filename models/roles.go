package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type RoleDB interface {
	GetAllRoles(pagination Pagination) ([]*Role, error)
	Count() (int, error)
}

type RoleService interface {
	RoleDB
}

var roleGoorm roleGorm

func newRoleGorm(db *gorm.DB) (*roleGorm, error) {
	roleGoorm.db = db
	return &roleGorm{
		db: db,
	}, nil
}

func NewRoleService(gD *gorm.DB) RoleDB {
	ug, err := newRoleGorm(gD)
	if err != nil {
		return nil
	}
	return ug
}


var _ RoleDB = &roleGorm{}

type roleGorm struct {
	db *gorm.DB
}

func (ug *roleGorm) GetAllRoles(pagination Pagination) ([]*Role, error) {
	offset := (pagination.Page - 1) * pagination.Limit
	var roles []*Role
	err := ug.db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil

}

// BASIC FUNCTIONS
func (role *Role) Create() error {
	return roleGoorm.db.Create(role).Error
}

func (role *Role) Delete() error {
	return roleGoorm.db.Delete(role).Error
}

func (role *Role) Update() error {
	return roleGoorm.db.Save(role).Error
}

func (role *Role) ModifyBalance(type_transaction string, ammount decimal.Decimal) error {
	return roleGoorm.db.Model(&Role{}).Where("id = ?", role.ID).Update(type_transaction, ammount).Error
}

// SEARCH BY
func (role *Role) ByID() (*Role, error) {
	db := roleGoorm.db.Where("id = ?", role.ID).First(role)
	err := first(db, role)
	return role, err
}

func (ug *roleGorm) Count() (int, error) {
	var roles int64
	err := ug.db.Table("roles").Count(&roles).Error
	return int(roles), err
}

type Role struct {
	NeftModel
	RoleName    string `gorm:"not null;unique_index" json:"rolename"`
	CanBanUsers bool   `gorm:"not null;default: false" json:"-"`
}
