package models

import (
	"quree/internal/models/enums"
)

type UUID string

type User struct {
	ChatID      string         `json:"chat_id"`
	PhoneNumber string         `json:"phone_number"`
	Role        enums.UserRole `json:"role"`
	QRCode      UUID           `json:"qr_code"`
}

// convert User to str

func (u *User) String() string {
	return u.ChatID
}

type Message struct {
	Content string            `json:"content"`
	Image   UUID              `json:"image"`
	Type    enums.MessageType `json:"type"`
	Sort    int32             `json:"sort"`
}

type UserEventVisit struct {
	UserID  UUID            `json:"user_id"`
	AdminID UUID            `json:"admin_id,omitempty"`
	QuizID  UUID            `json:"quiz_id,omitempty"`
	Type    enums.EventType `json:"type"`
}

type File struct {
	ID       UUID   `json:"id"`
	Filename string `json:"filename"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}
