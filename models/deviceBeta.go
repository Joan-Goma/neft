package models

import (
	"regexp"

	engine "github.com/JoanGTSQ/api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type DeviceDB interface {
	ByMac(mac string) (*Device, error)
	Create(device *Device) error
}

type DeviceService interface {
	DeviceDB
}

func newDeviceGorm(db *gorm.DB) (*deviceGorm, error) {
	return &deviceGorm{
		db: db,
	}, nil
}
func NewDeviceService(gD *gorm.DB) DeviceService {
	ug, err := newDeviceGorm(gD)
	if err != nil {
		return nil
	}
	hmac := engine.NewHMAC(hmacScretKey)
	uv := newDeviceValidator(ug, hmac)
	return &deviceService{
		DeviceDB: uv,
	}
}

type deviceService struct {
	DeviceDB
}
type deviceValidator struct {
	DeviceDB
	hmac       engine.HMAC
	emailRegex *regexp.Regexp
}

func newDeviceValidator(udb DeviceDB, hmac engine.HMAC) *deviceValidator {
	return &deviceValidator{
		DeviceDB:   udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

var _ DeviceDB = &deviceGorm{}

type deviceGorm struct {
	db *gorm.DB
}

func (tg *deviceGorm) Create(device *Device) error {
	err := tg.db.Create(device).Error
	if err != nil {
		return err
	}
	return nil
}

func (tg *deviceGorm) ByMac(mac string) (*Device, error) {
	var device Device
	db := tg.db.Where("mac = ?", string(mac)).First(&device)
	err := first(db, &device)
	return &device, err
}

type Device struct {
	NeftModel
	Name      string `gorm:"not null" json:"username"`
	Mac       string `gorm:"not null;unique_index" json:"mac"`
	Activated bool   `gorm:"default false"`
}
