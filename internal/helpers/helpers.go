package helpers

import (
	"quree/internal/models"
	"quree/internal/models/enums"
)

//write a function to convert models.UUID to string and return pointer

func UUIDToString(uuid models.UUID) *string {
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

// write a function to convert from string to enums.EventType by switching value between const of type EventType
// assume all values capitalized, and const names too
// if not found return only default value

func StringToEventType(str string) enums.EventType {
	switch str {
	case "EVENT":
		return enums.EVENT
	case "QUIZ":
		return enums.QUIZ
	default:
		return enums.EVENT
	}
}
