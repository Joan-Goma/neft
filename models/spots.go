package models

import (
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func (spot *Spot) Create() error {
	err := DBCONNECTION.Create(spot).Error
	if err != nil {
		return err
	}
	return nil
}

func (spot *Spot) Delete() error {
	return DBCONNECTION.Delete(&spot).Error
}

func (spot *Spot) Update() error {
	return DBCONNECTION.Save(spot).Error
}

// SEARCH BY ID

func (spot *Spot) ByID() error {
	if err := DBCONNECTION.
		Table("spots").
		Where("id = ?", spot.ID).
		First(spot).
		Error; err != nil {
		return err
	}
	return nil
}

func (spot *Spot) SearchNear() error {
	//TODO make search spots near a location
	return nil
}

func (ms *MultipleSpots) AllSpots() error {
	offset := (ms.Pagination.Page - 1) * ms.Pagination.Limit
	err := DBCONNECTION.Table("spots").
		Offset(offset).
		Limit(ms.Pagination.Limit).
		Order(ms.Pagination.Sort).
		Find(&ms.Spots).
		Error
	return err
}

func (ms *MultipleSpots) Count() error {
	err := DBCONNECTION.Table("spots").Count(&ms.Quantity).Error
	return err
}

type MultipleSpots struct {
	Pagination Pagination
	Spots      []*Spot
	Quantity   int64
}

type Spot struct {
	NeftModel
	Description string          `json:"description"`
	Latitude    decimal.Decimal `json:"latitude"  sql:"type:decimal(9,8);"`
	Longitude   decimal.Decimal `json:"longitude" sql:"type:decimal(9,8);"`
	Visible     bool            `gorm:"default: true" json:"visible"`
}
