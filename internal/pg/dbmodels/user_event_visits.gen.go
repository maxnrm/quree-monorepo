// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbmodels

import (
	"time"
)

const TableNameUserEventVisit = "user_event_visits"

// UserEventVisit mapped from table <user_event_visits>
type UserEventVisit struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	DateCreated time.Time `gorm:"column:date_created;type:timestamp with time zone;not null" json:"date_created"`
	QuizID      *string   `gorm:"column:quiz_id;type:character varying(255)" json:"quiz_id"`
	EventType   *string   `gorm:"column:event_type;type:character varying(255)" json:"event_type"`
	UserID      *string   `gorm:"column:user_id;type:uuid" json:"user_id"`
	AdminID     *string   `gorm:"column:admin_id;type:uuid" json:"admin_id"`
}

// TableName UserEventVisit's table name
func (*UserEventVisit) TableName() string {
	return TableNameUserEventVisit
}