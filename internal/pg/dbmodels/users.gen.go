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
	ID               string     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	DateCreated      time.Time  `gorm:"column:date_created;type:timestamp with time zone;not null" json:"date_created"`
	QrCode           string     `gorm:"column:qr_code;type:uuid;not null" json:"qr_code"`
	ChatID           string     `gorm:"column:chat_id;type:character varying(255);not null;default:NULL" json:"chat_id"`
	QuizCityName     *string    `gorm:"column:quiz_city_name;type:character varying(255)" json:"quiz_city_name"`
	IsFinished       bool       `gorm:"column:is_finished;type:boolean;not null" json:"is_finished"`
	DateQuizFinished *time.Time `gorm:"column:date_quiz_finished;type:timestamp without time zone" json:"date_quiz_finished"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
