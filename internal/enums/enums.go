package enums

type UserType string
type MessageType string
type ImageType string

const (
	PROFILE ImageType = "PROFILE"
	QR      ImageType = "QR"
	ATTACH  ImageType = "ATTACH"
)

const (
	ADMIN UserType = "ADMIN"
	USER  UserType = "USER"
)

const (
	HELP       MessageType = "HELP"
	LORE_EVENT MessageType = "LORE_EVENT"
)
