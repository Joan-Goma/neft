package models

import (
	engine "github.com/JoanGTSQ/api"
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

func (message *UserMessage) RegisterMessage() {
	err := DBCONNECTION.Create(&message).Error
	if err != nil {
		engine.Error.Fatalln(err)
	}
}
