package enums

// create idiomatic golang enums for User.Role
// values could be ADMIN and USER
type UserRole string
type MessageType string
type EventType string

const (
	ADMIN UserRole = "ADMIN"
	USER  UserRole = "USER"
)

const (
	HELP             MessageType = "HELP"
	START            MessageType = "START"
	LORE_EVENT1      MessageType = "LORE_EVENT1"
	LORE_EVENT2      MessageType = "LORE_EVENT2"
	LORE_EVENT3      MessageType = "LORE_EVENT3"
	LORE_EVENT4      MessageType = "LORE_EVENT4"
	LORE_EVENT_EXTRA MessageType = "LORE_EVENT_EXTRA"
	LORE_EVENT_QUIZ  MessageType = "LORE_EVENT_QUIZ"
)

const (
	EVENT EventType = "EVENT"
	QUIZ  EventType = "QUIZ"
)
