package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"library/internal/config"
)

func main() {
	cfg := config.MustLoad()
	var migrationsPath = cfg.Postgres.MigrationsPath
	if migrationsPath == "" {
		panic("MIGRATIONS_PATH is required")
	}
	m, err := migrate.New("file://"+migrationsPath, cfg.Postgres.StorageURL)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
