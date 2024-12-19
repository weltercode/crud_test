package app

import (
	"crud_test/internal/database/postgres"
	"crud_test/internal/repositories"
	"crud_test/internal/transport/rest"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	config   *Config
	router   *mux.Router
	handler  *rest.Handler
	db       *sql.DB
	taskRepo repositories.TaskRepositoryInterface
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
	taskRepo := repositories.NewTaskRepository(db)

	return &App{
		config:   config,
		router:   router,
		handler:  rest.NewHandler(router, taskRepo),
		db:       db,
		taskRepo: taskRepo,
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
	app.router.HandleFunc("/task/delete/{id:[0-9]+}", app.handler.DeleteTaskHandler).Methods("GET").Name("task_delete")
	app.router.HandleFunc("/task/start/{id:[0-9]+}", app.handler.StartTask).Methods("GET").Name("task_start")
	app.router.HandleFunc("/task/end/{id:[0-9]+}", app.handler.EndTask).Methods("GET").Name("task_end")
	http.Handle("/", app.router)

	srv := &http.Server{
		Addr:    ":" + app.config.AppPort,
		Handler: app.router,
	}

	log.Printf("SERVER STARTED localhost:%s AT %s", app.config.AppPort, time.Now().Format(time.RFC3339))

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

func (app *App) CreateTables() {
	defer app.Shutdown()
	query := `
    CREATE TABLE IF NOT EXISTS"tasks" (
        "id" SERIAL NOT NULL PRIMARY KEY,
        "title" VARCHAR(255) NOT NULL,
        "description" VARCHAR(255) NULL,
        "starttime" TIMESTAMP NULL,
        "endtime" TIMESTAMP NULL
    );`

	_, err := app.db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	} else {
		log.Println("Tables created successfully!")
	}
}
