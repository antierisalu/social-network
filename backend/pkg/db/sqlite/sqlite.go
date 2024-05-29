package sqlite

import (
    "log"
    migrate "github.com/golang-migrate/migrate/v4"
    // _ "github.com/golang-migrate/migrate/v4/database/sqlite3"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	//TODO START DB AND EVERYTHING ELSE
	m, err := migrate.New(
        // "file://path/to/migrations",
        // "sqlite3://path/to/your/database.db",
		"../../db/migrations/",
		"./",
    )
    if err != nil {
        log.Fatalf("Migration init error: %v\n", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Migration up error: %v\n", err)
    }

    log.Println("Migration completed successfully")
}
