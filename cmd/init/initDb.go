package main

import (
	"crud_test/internal/app"
	"log"
	"os"

	"github.com/caarlos0/env/v11"

	"github.com/joho/godotenv"
)

func main() {
	working_env := os.Getenv("APP_ENV")
	if working_env == "production" {
		if err := godotenv.Load(".env.production"); err != nil {
			log.Println("No .env.production file found, using system environment variables")
		}
	} else {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	// Parse configuration from environment variables
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Initialize the app and create tables
	app := app.CreateApp(&cfg)
	log.Println("Initializing the database...")
	app.CreateTables()
	log.Println("Database initialization complete.")
}
