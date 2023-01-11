package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type SpotDB interface {
	AllSpots(pagination Pagination, userID uint) ([]*Spot, error)
	Count() (int, error)
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

func (ug *spotGorm) AllSpots(pagination Pagination, userID uint) ([]*Spot, error) {
	var spot []*Spot
	offset := (pagination.Page - 1) * pagination.Limit
	err := ug.db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort).Find(&spot).Error
	return spot, err
}

func (tg *spotGorm) Count() (int, error) {
	var spots int64
	err := tg.db.Table("spots").Count(&spots).Error
	return int(spots), err
}

type Spot struct {
	NeftModel
	Description string `json:"description"`
	Latitude    int64  `json:"latitude"`
	Longitude   int64  `json:"longitude"`
	Visible     bool   `gorm:"default: true" json:"visible"`
}
