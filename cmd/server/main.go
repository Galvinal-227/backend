package main

import (
	"log"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/database"
	"pkl-attendance/backend/routes"
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

	r := routes.Setup(cfg, db)
	log.Printf("server listening on :%s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("run: %v", err)
	}
}