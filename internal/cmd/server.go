package main

import (
	"crud_test/internal/app"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// file:///C:/Users/welte/Downloads/practice+%231.pdf
func main() {
	// Load .env file into environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println(cfg)
	application := app.CreateApp(&cfg)
	application.Run()
}
