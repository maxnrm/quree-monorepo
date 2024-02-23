package main

import (
	"fmt"
	"quree/internal/pg"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type test struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func main() {

	// create relation test
	pg.DB.AutoMigrate(&test{})

	testModel := &test{
		Name: "test",
	}

	pg.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&testModel)

	fmt.Println()

}
