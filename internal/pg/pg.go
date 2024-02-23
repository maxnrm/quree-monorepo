package pg

import (
	"fmt"
	"quree/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	DSN string
}

var pgConfig = &dbConfig{
	DSN: config.POSTGRES_CONN_STRING,
}

var DB = Init(pgConfig)

func Init(c *dbConfig) *gorm.DB {

	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprint("Failed to connect to database at dsn: ", c.DSN))
	}

	return db
}
