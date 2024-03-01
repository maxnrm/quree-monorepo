package enums

// create idiomatic golang enums for User.Role
// values could be ADMIN and USER
type UserRole string

const (
	ADMIN UserRole = "ADMIN"
	USER  UserRole = "USER"
)
