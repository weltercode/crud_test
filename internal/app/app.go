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
	config  *Config
	router  *mux.Router
	handler *rest.Handler
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
	router := mux.NewRouter()
	return &App{
		config:  config,
		router:  router,
		handler: rest.NewHandler(router),
	}
}
func (app *App) Run() {
	log.Println("App Running")

	app.router.HandleFunc("/", app.handler.HomeHandler).Name("home")
	app.router.HandleFunc("/tasks", app.handler.TasksListHandler).Methods("GET").Name("tasks_list")
	app.router.HandleFunc("/task/{id:[0-9]+}", app.handler.TaskViewHandler).Methods("GET").Name("task_view")
	app.router.HandleFunc("/task/", app.handler.TaskViewHandlerzz).Methods("GET").Name("task_view_zz")
	app.router.HandleFunc("/task/new", app.handler.TaskViewHandler).Methods("GET").Name("task_new")
	app.router.HandleFunc("/task/new", app.handler.TaskViewHandler).Methods("GET").Name("task_new")
	app.router.HandleFunc("/task/add", app.handler.TaskViewHandler).Methods("GET").Name("task_add")
	app.router.HandleFunc("/task/login", app.handler.LoginHandler).Methods("GET").Name("login")
	http.Handle("/", app.router)

	srv := &http.Server{
		Addr:    ":" + app.config.AppPort,
		Handler: app.router,
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
