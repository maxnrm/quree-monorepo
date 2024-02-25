package enums

type UserType string
type MessageType string

const (
	ADMIN UserType = "ADMIN"
	USER  UserType = "USER"
)

const (
	HELP       MessageType = "HELP"
	LORE_EVENT MessageType = "LORE_EVENT"
)
