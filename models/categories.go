package models

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"regexp"
)

type CategoryDB interface {
	ByID(id uint) (*Category, error)
	AllCategories(pagination Pagination) ([]*Category, error)

	Create(category *Category) error
	Update(category *Category) error
	Delete(category *Category) error
	Count() (int, error)
}

type CategoryService interface {
	CategoryDB
}

func newCategoryGorm(db *gorm.DB) (*categoryGorm, error) {
	return &categoryGorm{
		db: db,
	}, nil
}
func NewCategoryService(gD *gorm.DB) CategoryService {
	ug, err := newCategoryGorm(gD)
	if err != nil {
		return nil
	}
	hmac := engine.NewHMAC(hmacScretKey)
	uv := newCategoryValidator(ug, hmac)
	return &categoryService{
		CategoryDB: uv,
	}
}

type categoryService struct {
	CategoryDB
}
type categoryValidator struct {
	CategoryDB
	hmac       engine.HMAC
	emailRegex *regexp.Regexp
}

func newCategoryValidator(udb CategoryDB, hmac engine.HMAC) *categoryValidator {
	return &categoryValidator{
		CategoryDB: udb,
		hmac:       hmac,
	}
}

var _ CategoryDB = &categoryGorm{}

type categoryGorm struct {
	db *gorm.DB
}

func (tg *categoryGorm) Create(category *Category) error {
	err := tg.db.Create(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (tg *categoryGorm) Delete(category *Category) error {
	return tg.db.Delete(&category).Error
}

func (tg *categoryGorm) Update(category *Category) error {
	return tg.db.Save(category).Error
}

// SEARCH BY ID
func (ug *categoryGorm) ByID(id uint) (*Category, error) {
	var category Category
	err := ug.db.Where("id = ?", id).First(&category).Error
	return &category, err
}

func (ug *categoryGorm) AllCategories(pagination Pagination) ([]*Category, error) {
	var category []*Category
	offset := (pagination.Page - 1) * pagination.Limit
	err := ug.db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort).Find(&category).Error
	return category, err
}

func (tg *categoryGorm) Count() (int, error) {
	var categories int64
	err := tg.db.Table("category").Count(&categories).Error
	return int(categories), err
}

type Category struct {
	NeftModel
	Name        string `json:"name"`
	Description string `json:"description"`
}
