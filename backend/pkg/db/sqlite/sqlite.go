package sqlite

import (
	"database/sql"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

func DatabaseInit() {
	m, err := migrate.New(
		"file://pkg/db/migrations/sqlite",
		"sqlite3://pkg/db/sqlite/database.db",
	)
	if err != nil {
		log.Fatalf("Migration init error: %v\n", err)
	}
	//sqlite cli support is gone so just put the number of the db version u wanna use
	if err := m.Migrate(5); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration up error: %v\n", err)
	}
	DB, err = sql.Open("sqlite3", "pkg/db/sqlite/database.db")
	if err != nil {
		log.Println("error validating sql.Open arguments")
		panic(err.Error())
	}

	log.Println("Migration completed successfully")
}
