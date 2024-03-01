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

	"github.com/google/uuid"

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
func (pg *pg) CreateUser(user *models.User) error {

	result := pg.Create(&dbmodels.User{
		ID:          uuid.New().String(),
		ChatID:      user.ChatID,
		PhoneNumber: &user.PhoneNumber,
		DateCreated: time.Now(),
		Role:        string(user.Role),
		QrCode:      string("f5a879d4-147d-4740-a4e8-b8bb8f5a791c"),
	})

	return result.Error
}

// function to get user from db using ChatID, transform in models.User struct and return

func (pg *pg) GetUserByChatID(chatID string) *models.User {

	var user dbmodels.User
	result := pg.Where("chat_id = ?", chatID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var pn string
	if user.PhoneNumber != nil {
		pn = *user.PhoneNumber
	} else {
		pn = ""
	}

	return &models.User{
		ChatID:      user.ChatID,
		PhoneNumber: pn,
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
