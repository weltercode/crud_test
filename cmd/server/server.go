package main

import (
	"crud_test/internal/app"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// file:///C:/Users/welte/Downloads/practice+%231.pdf
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
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	application := app.CreateApp(&cfg)
	application.Run()
}
