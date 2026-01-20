package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ParkPawapon/mhp-be/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	path := flag.String("path", "migrations", "path to migrations")
	action := flag.String("action", "up", "up | down | force")
	steps := flag.Int("steps", 0, "number of steps for up/down or version for force")
	flag.Parse()

	m, err := migrate.New("file://"+*path, cfg.DB.URL())
	if err != nil {
		log.Fatalf("migrate init failed: %v", err)
	}

	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "force":
		if *steps == 0 {
			fmt.Println("force requires -steps version")
			os.Exit(1)
		}
		err = m.Force(*steps)
	default:
		fmt.Println("unknown action:", *action)
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
}
