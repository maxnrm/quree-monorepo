package pg

import (
	"fmt"
	"quree/config"
	"quree/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB = Init(config.POSTGRES_CONN_STRING)

func Init(connString string) *gorm.DB {

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprint("Failed to connect to database at dsn: ", connString))
	}

	return db
}

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(&models.User{}, &models.Message{}, &models.UserEventVisit{})
}
