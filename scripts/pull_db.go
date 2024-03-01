// configuration.go
package main

import (
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const dsn = "postgres://quree:qureequree@127.0.0.1:5432/quree"

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/pg/dbquery", // output directory, default value is ./query
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		ModelPkgPath:      "dbmodels",
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connection err: %v\n", err)
	}

	g.UseDB(db)

	g.WithTableNameStrategy(func(tableName string) (targetTableName string) {
		if strings.HasPrefix(tableName, "_") {
			return ""
		}
		return tableName
	})

	g.ApplyBasic(
		g.GenerateModelAs("users", "User"),
		g.GenerateModelAs("messages", "Message"),
		g.GenerateModelAs("user_event_visits", "UserEventVisit"),
		g.GenerateModelAs("directus_files", "File"),
	)

	// Execute the generator
	g.Execute()
}
