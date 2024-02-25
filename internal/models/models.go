package models

import (
	"quree/internal/enums"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Role         enums.UserType `sql:"type:enum('ADMIN', 'USER');'default:USER'"`
	ChatID       int64          `gorm:"unique"`
	ProfileQR    string
	ProfileImage string
}

type Message struct {
	gorm.Model
	Type       enums.MessageType `sql:"type:enum('HELP', 'LORE_EVENT');default:'LORE_EVENT'"`
	Content    string
	Attachment string
	Order      int
}

type UserEventVisit struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	AdminID uint
	Admin   User `gorm:"foreignKey:AdminID"`
}
