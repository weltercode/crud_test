package app

import (
	"crud_test/internal/database/postgres"
	"crud_test/internal/transport/rest"
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
	log.Println("App Running")

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rest.HomeHandler(w, r, router) // Pass the router itself
	})
	router.HandleFunc("/tasks", rest.TasksListHandler).Methods("GET").Name("tasks_list")
	router.HandleFunc("/task/{id:[0-9]+}", rest.TaskViewHandler).Methods("GET").Name("task_view")
	router.HandleFunc("/task/new", rest.TaskViewHandler).Methods("GET").Name("task_new")
	router.HandleFunc("/task/add", rest.TaskViewHandler).Methods("GET").Name("task_add")
	http.Handle("/", router)

	srv := &http.Server{
		Addr:    ":" + app.config.AppPort,
		Handler: router,
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	db := postgres.NewDbConnect(postgres.ConnectionConfig{
		Host:   app.config.DbIp,
		Port:   app.config.DbPort,
		DbName: app.config.DbName,
		User:   app.config.DbUser,
		Pass:   app.config.DbPass,
	})
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connected")
	}
}
