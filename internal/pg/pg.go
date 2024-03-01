package pg

import (
	"errors"
	"fmt"
	"quree/config"
	"time"

	"quree/internal/helpers"
	"quree/internal/models"

	"quree/internal/pg/dbmodels"
	"quree/internal/pg/dbquery"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (pg *pg) CreateUser(user *models.User) error {

	result := pg.Clauses(clause.OnConflict{DoNothing: true}).Create(&dbmodels.User{
		ChatID:      user.ChatID,
		PhoneNumber: &user.PhoneNumber,
		DateCreated: time.Now(),
		Role:        string(user.Role),
		QrCode:      string(user.QRCode),
	})

	return result.Error
}

// function to get user from db using ChatID, transform in models.User struct and return

func (pg *pg) GetUserByChatID(chatID int64) *models.User {

	var user dbmodels.User
	result := pg.Where("chat_id = ?", chatID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &models.User{
		ChatID:      user.ChatID,
		PhoneNumber: *user.PhoneNumber,
		Role:        helpers.StringToUserRole(user.Role),
		QRCode:      models.UUID(user.QrCode),
	}
}

func (pg *pg) Close() {
	sqlDB, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
