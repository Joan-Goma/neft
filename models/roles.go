package models

import (
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func (mr *MultipleRoles) GetAllRoles() error {
	offset := (mr.Pagination.Page - 1) * mr.Pagination.Limit
	err := DBCONNECTION.
		Table("roles").
		Offset(offset).
		Limit(mr.Pagination.Limit).
		Order(mr.Pagination.Sort).
		Find(&mr.Roles).
		Error
	if err != nil {
		return err
	}
	return nil

}

// BASIC FUNCTIONS
func (role *Role) Create() error {
	return DBCONNECTION.Create(role).Error
}

func (role *Role) Delete() error {
	return DBCONNECTION.Delete(role).Error
}

func (role *Role) Update() error {
	return DBCONNECTION.Save(role).Error
}

func (role *Role) ModifyBalance(typeTransaction string, ammount decimal.Decimal) error {
	return DBCONNECTION.Model(&Role{}).
		Where("id = ?", role.ID).
		Update(typeTransaction, ammount).
		Error
}

// SEARCH BY
func (role *Role) ByID() (*Role, error) {
	db := DBCONNECTION.
		Table("roles").
		Where("id = ?", role.ID).
		First(role)
	err := first(db, role)
	return role, err
}

func (mr *MultipleRoles) Count() error {
	err := DBCONNECTION.
		Table("roles").
		Count(&mr.Roles).
		Error
	return err
}

type MultipleRoles struct {
	Pagination Pagination
	Roles      []*Role
	Quantity   int64
}
type Role struct {
	NeftModel
	RoleName    string `gorm:"not null;unique_index" json:"rolename"`
	CanBanUsers bool   `gorm:"not null;default: false" json:"-"`
}
