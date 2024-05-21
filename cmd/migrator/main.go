package main

import (
	"WalletRieltaTestTask/config"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {

	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "./schema", "path to migrations")
	flag.Parse()

	if migrationsPath == "" {
		panic("Migrations-path is required")
	}

	cfg := config.MustLoad()

	if err := godotenv.Load(); err != nil {
		panic("Failed reading db password")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf(cfg.PG.URL),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("Migrations applied")
}
