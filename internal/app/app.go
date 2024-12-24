package app

import (
	"crud_test/internal/database/postgres"
	"crud_test/internal/logger"
	"crud_test/internal/repositories"
	"crud_test/internal/transport/rest"
	"fmt"

	"database/sql"
	"net/http"
	"time"

	_ "crud_test/docs" // which is the generated folder after swag init

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	config   *Config
	router   *mux.Router
	handler  *rest.Handler
	db       *sql.DB
	taskRepo repositories.TaskRepositoryInterface
	logger   logger.LoggerInterface
}
type Config struct {
	DbIp    string `env:"DATABASE_IP"`
	DbPort  string `env:"DATABASE_PORT"`
	DbName  string `env:"DATABASE_NAME"`
	DbUser  string `env:"DATABASE_USER"`
	DbPass  string `env:"DATABASE_PASS"`
	AppPort string `env:"APP_PORT"`
}

func CreateApp(config *Config) *App {
	router := mux.NewRouter()
	logger := logger.NewSlogLogger()
	db := postgres.NewDbConnect(postgres.ConnectionConfig{
		Host:   config.DbIp,
		Port:   config.DbPort,
		DbName: config.DbName,
		User:   config.DbUser,
		Pass:   config.DbPass,
	}, logger)

	if err := db.Ping(); err != nil {
		logger.Error("Fail to connect DB", err)
	} else {
		logger.Info("Database connected", err)
	}
	taskRepo := repositories.NewTaskRepository(db, logger)

	return &App{
		config:   config,
		router:   router,
		handler:  rest.NewHandler(router, taskRepo, logger),
		db:       db,
		taskRepo: taskRepo,
		logger:   logger,
	}
}
func (app *App) Run() {
	defer app.Shutdown()
	app.logger.Info("App Running")

	app.router.HandleFunc("/", app.handler.HomeHandler).Name("home")
	app.router.HandleFunc("/tasks", app.handler.TasksListHandler).Methods("GET").Name("tasks_list")
	app.router.HandleFunc("/task/{id:[0-9]+}", app.handler.TaskViewHandler).Methods("GET").Name("task_view")
	app.router.HandleFunc("/task/new", app.handler.TaskViewHandler).Methods("GET").Name("task_new")
	app.router.HandleFunc("/task/save", app.handler.TaskSaveHandler).Methods("POST").Name("task_save")
	app.router.HandleFunc("/task/login", app.handler.LoginHandler).Methods("GET").Name("login")
	app.router.HandleFunc("/task/delete/{id:[0-9]+}", app.handler.DeleteTaskHandler).Methods("GET").Name("task_delete")
	app.router.HandleFunc("/task/start/{id:[0-9]+}", app.handler.StartTask).Methods("GET").Name("task_start")
	app.router.HandleFunc("/task/end/{id:[0-9]+}", app.handler.EndTask).Methods("GET").Name("task_end")
	app.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app.router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	http.Handle("/", app.router)

	srv := &http.Server{
		Addr:    ":" + app.config.AppPort,
		Handler: app.router,
	}

	app.logger.Info(fmt.Sprintf("SERVER STARTED localhost:%s AT %s", app.config.AppPort, time.Now().Format(time.RFC3339)))

	if err := srv.ListenAndServe(); err != nil {
		app.logger.Error("Server fail to start", err)
	}

}
func (app *App) Shutdown() {
	if app.db != nil {
		app.logger.Info("Closing database connection...")
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
		app.logger.Error("Failed to create tables", err)
	} else {
		app.logger.Info("Tables created successfully!")
	}
}
