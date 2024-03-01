package pg

import (
	"fmt"
	"quree/config"

	"quree/internal/models/enums"
	"quree/internal/pg/dbmodels"
	"quree/internal/pg/dbquery"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pg struct {
	*gorm.DB
	q *dbquery.Query
}

var DB = Init(config.POSTGRES_CONN_STRING)

func Init(connString string) *pg {

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprint("Failed to connect to database at dsn: ", connString))
	}

	var pg *pg = &pg{
		DB: db,
		q:  dbquery.Use(db),
	}

	return pg
}

// function to create user
func (pg *pg) CreateUser(chatID string, phone_number *string, role enums.UserRole, profilePic *string, qrCode *string) {

	roleStr := string(role)

	user := dbmodels.User{
		ChatID:      &chatID,
		PhoneNumber: phone_number,
		Role:        &roleStr,
	}

	pg.q.User.Create(&user)
}

func (pg *pg) Close() {
	sqlDB, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
