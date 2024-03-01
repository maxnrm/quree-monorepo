package pg

import (
	"fmt"
	"quree/config"
	"time"

	"quree/internal/models"

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
func (pg *pg) CreateUser(user *models.User) {

	pg.q.User.Create(&dbmodels.User{
		ChatID:      user.ChatID,
		PhoneNumber: &user.PhoneNumber,
		DateCreated: time.Now(),
		Role:        string(user.Role),
		QrCode:      string(user.QRCode),
	})

}

func (pg *pg) Close() {
	sqlDB, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
