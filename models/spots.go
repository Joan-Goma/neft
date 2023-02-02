package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type SpotDB interface {
}

type SpotService interface {
	SpotDB
}

var gormSpot spotGorm

func newSpotGorm(db *gorm.DB) (*spotGorm, error) {
	gormSpot.db = db
	return &spotGorm{
		db: db,
	}, nil
}
func NewSpotService(gD *gorm.DB) SpotService {
	ug, err := newSpotGorm(gD)
	if err != nil {
		return nil
	}
	return &spotService{
		SpotDB: ug,
	}
}

type spotService struct {
	SpotDB
}

var _ SpotDB = &spotGorm{}

type spotGorm struct {
	db *gorm.DB
}

func (spot *Spot) Create() error {
	err := gormSpot.db.Create(spot).Error
	if err != nil {
		return err
	}
	return nil
}

func (spot *Spot) Delete() error {
	return gormSpot.db.Delete(&spot).Error
}

func (spot *Spot) Update() error {
	return gormSpot.db.Save(spot).Error
}

// SEARCH BY ID

func (spot *Spot) ByID() error {
	if err := gormSpot.db.Where("id = ?", spot.ID).First(spot).Error; err != nil {
		return err
	}
	return nil
}

func (spot *Spot) SearchNear(offset uint) error {
	//TODO make search spots near a location
	return nil
}

func (ms *MultipleSpots) AllSpots(userID uint) error {
	offset := (ms.Pagination.Page - 1) * ms.Pagination.Limit
	err := gormSpot.db.Offset(offset).Limit(ms.Pagination.Limit).Order(ms.Pagination.Sort).Find(&ms.Spots).Error
	return err
}

func (ms *MultipleSpots) Count() error {
	err := gormSpot.db.Table("spots").Count(&ms.Quantity).Error
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
