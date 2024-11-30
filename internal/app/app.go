package app

import (
	"crud_test/internal/transport/rest"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	config *Config
}
type Config struct {
	DbIp    string `env:"DATABASE_IP"`
	DbPort  string `env:"DATABASE_PORT"`
	DbName  string `env:"DATABASE_NAME"`
	DbUser  string `env:"DATABASE_USER"`
	DbPass  string `env:"DATABASE_PASS"`
	AppPort string `env:"APP_PORT" envDefault:8080`
}

func CreateApp(config *Config) *App {
	return &App{config: config}
}
func (app *App) Run() {
	fmt.Println("Running")

	r := mux.NewRouter()
	r.HandleFunc("/", rest.HomeHandler)
	r.HandleFunc("/products", rest.ProductsHandler)
	r.HandleFunc("/articles", rest.ArticlesHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr:    ":" + app.config.AppPort,
		Handler: r,
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
