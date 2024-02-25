package main

import "quree/internal/pg"

func main() {
	db := pg.DB
	pg.Migrate(db)
}
