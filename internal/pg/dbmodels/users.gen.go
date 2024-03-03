// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbmodels

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	DateCreated time.Time `gorm:"column:date_created;type:timestamp with time zone;not null" json:"date_created"`
	QrCode      *string   `gorm:"column:qr_code;type:uuid" json:"qr_code"`
	Role        string    `gorm:"column:role;type:character varying(255);not null;default:USER" json:"role"`
	ChatID      *string   `gorm:"column:chat_id;type:character varying(255);default:NULL" json:"chat_id"`
	PhoneNumber *string   `gorm:"column:phone_number;type:character varying(255);default:NULL" json:"phone_number"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
