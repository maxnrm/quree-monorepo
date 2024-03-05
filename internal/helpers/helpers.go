package helpers

import (
	"log"
	"quree/internal/models"
	"quree/internal/models/enums"

	tele "gopkg.in/telebot.v3"
)

//write a function to convert models.UUID to string and return pointer

func UUIDToString(uuid models.UUID) *string {
	if uuid == "" {
		return nil
	}

	str := string(uuid)
	return &str
}

// write a function to convert from string to enums.UserRole by switching value between const of type UserRole
// assume all values capitalized, and const names too
// if not found return only default value

func StringToUserRole(str string) enums.UserRole {
	switch str {
	case "ADMIN":
		return enums.ADMIN
	case "USER":
		return enums.USER
	default:
		return enums.USER
	}
}

// write a function to convert from string to enums.MessageType by switching value between const of type MessageType
// assume all values capitalized, and const names too
// if not found return only default value

func StringToMessageType(str string) enums.MessageType {
	switch str {
	case "HELP":
		return enums.HELP
	case "START":
		return enums.START
	case "LORE_EVENT1":
		return enums.LORE_EVENT1
	case "LORE_EVENT2":
		return enums.LORE_EVENT2
	case "LORE_EVENT3":
		return enums.LORE_EVENT3
	case "LORE_EVENT4":
		return enums.LORE_EVENT4
	case "LORE_EVENT_EXTRA":
		return enums.LORE_EVENT_EXTRA
	case "LORE_EVENT_QUIZ":
		return enums.LORE_EVENT_QUIZ
	default:
		return enums.HELP
	}
}

// function to get models.MessageType by count of UserEventVisit
// 1 - LORE_EVENT1, 2 - LORE_EVENT2, 3 - LORE_EVENT3, 4 - LORE_EVENT4, 5 and more - LORE_EVENT_EXTRA

func GetMessageTypeByCount(count int) enums.MessageType {
	switch {
	case count == 1:
		return enums.LORE_EVENT1
	case count == 2:
		return enums.LORE_EVENT2
	case count == 3:
		return enums.LORE_EVENT3
	case count == 4:
		return enums.LORE_EVENT4
	case count > 4:
		return enums.LORE_EVENT_EXTRA
	default:
		return enums.START
	}
}

func BotMiniLogger() tele.MiddlewareFunc {
	l := log.Default()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			chatID := c.Chat().ID
			text := c.Message().Text
			l.Println(chatID, text, "ok")
			return next(c)
		}
	}
}
