package main

import (
	"crud_test/internal/app"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
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
