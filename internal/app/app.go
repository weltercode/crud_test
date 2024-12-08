package app

import (
	"crud_test/internal/database/postgres"
	"crud_test/internal/transport/rest"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	config  *Config
	router  *mux.Router
	handler *rest.Handler
	db      *sql.DB
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
	db := postgres.NewDbConnect(postgres.ConnectionConfig{
		Host:   config.DbIp,
		Port:   config.DbPort,
		DbName: config.DbName,
		User:   config.DbUser,
		Pass:   config.DbPass,
	})
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connected")
	}

	return &App{
		config:  config,
		router:  router,
		handler: rest.NewHandler(router, db),
		db:      db,
	}
}
func (app *App) Run() {
	defer app.Shutdown()
	log.Println("App Running")

	app.router.HandleFunc("/", app.handler.HomeHandler).Name("home")
	app.router.HandleFunc("/tasks", app.handler.TasksListHandler).Methods("GET").Name("tasks_list")
	app.router.HandleFunc("/task/{id:[0-9]+}", app.handler.TaskViewHandler).Methods("GET").Name("task_view")
	app.router.HandleFunc("/task/new", app.handler.TaskViewHandler).Methods("GET").Name("task_new")
	app.router.HandleFunc("/task/save", app.handler.TaskSaveHandler).Methods("POST").Name("task_save")
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
}
func (app *App) Shutdown() {
	if app.db != nil {
		log.Println("Closing database connection...")
		app.db.Close()
	}
}
