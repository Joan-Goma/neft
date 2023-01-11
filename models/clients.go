package models

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type UserMessage struct {
	SenderID uint      `gorm:"not null"`
	Type     string    `gorm:"not null"`
	Sender   User      `gorm:"foreignkey:SenderID"`
	Message  string    `gorm:"not null"`
	Receiver uuid.UUID `gorm:"not null"`
}

type MessageDB interface {
}

var gormMessage postGorm

type MessageService interface {
	MessageDB
}

var _ MessageDB = &messageGorm{}

type messageGorm struct {
	db *gorm.DB
}

func NewMessageService(gD *gorm.DB) MessageService {
	ug, err := newGormMessage(gD)
	if err != nil {
		return nil
	}
	return &messageService{
		MessageDB: ug,
	}
}

type messageService struct {
	MessageDB
}

func newGormMessage(db *gorm.DB) (*postGorm, error) {
	gormMessage.db = db
	return &postGorm{
		db: db,
	}, nil
}

func (message UserMessage) RegisterMessage() {
	err := gormMessage.db.Create(&message).Error
	if err != nil {
		engine.Error.Fatalln(err)
	}
}
