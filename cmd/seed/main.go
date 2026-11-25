package main

import (
	"log"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/database"
	"pkl-attendance/backend/seeders"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := seeders.SeedAdmin(db, cfg); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Println("seed done")
}