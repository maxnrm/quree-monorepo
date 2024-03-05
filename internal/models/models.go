package models

import (
	tele "gopkg.in/telebot.v3"
)

type UUID string

type User struct {
	ID           UUID   `json:"id"`
	ChatID       string `json:"chat_id"`
	QRCode       UUID   `json:"qr_code"`
	QuizCityName string `json:"quiz_city_name"`
}

type Admin struct {
	ID     UUID   `json:"id"`
	ChatID string `json:"chat_id"`
}

type Message struct {
	Text        *string           `json:"text"`
	Caption     *string           `json:"caption"`
	Photo       *tele.Photo       `json:"photo"`
	SendOptions *tele.SendOptions `json:"send_options"`
	Variant     int               `json:"variant"`
}

type UserEventVisit struct {
	DateCreated string `json:"date_created"`
	UserChatID  string `json:"user_id"`
	AdminChatID string `json:"admin_id"`
}

type File struct {
	ID       UUID   `json:"id"`
	Filename string `json:"filename"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type QRCodeMessage struct {
	UserChatID  string `json:"chat_id"`
	AdminChatID string `json:"admin_chat_id,omitempty"`
	QRCodeID    UUID   `json:"qr_code_id"`
}
