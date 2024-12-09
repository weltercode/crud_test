package rest

import (
	"crud_test/internal/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
	tmpl   *template.Template
	db     *sql.DB
}

func NewHandler(router *mux.Router, db *sql.DB) *Handler {
	tmpl := template.New("")
	return &Handler{
		router: router,
		tmpl:   tmpl,
		db:     db,
	}
}

// BaseHandler renders the base layout and dynamically includes specific content
func (h *Handler) BaseHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {

	w.WriteHeader(http.StatusOK)

	data["Home_page_url"] = h.getHrefByRouteName("home").String()
	data["Tasks_page_url"] = h.getHrefByRouteName("tasks_list").String()
	data["Sign_page_url"] = h.getHrefByRouteName("login").String()

	if err := h.tmpl.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HomeHandler handles the "/" route and injects specific content
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	h.tmpl, err = template.ParseFiles("templates/base.html", "templates/header.html", "templates/footer.html", "templates/home.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Title":       "Welcome to simple task tracker",
		"Description": "Manage your tasks efficiently",
		"newTaskHref": h.getHrefByRouteName("task_new").String(),
	}

	h.BaseHandler(w, r, data)
}

func formatTimestamp(ts int32) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}
func formatTime(t time.Time) string {
	return fmt.Sprintf("%02d.%02d.%04d %02d:%02d:%02d",
		t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second())
}

// TasksListHandler handles the "/tasks" route
func (h *Handler) TasksListHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	funcMap := template.FuncMap{
		"formatTimestamp": formatTimestamp,
		"formatT":         formatTime,
	}
	h.tmpl.Funcs(funcMap)
	h.tmpl, err = h.tmpl.ParseFiles(
		"templates/base.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/tasks_list.html",
	)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := []models.Task{
		{
			ID:          "1",
			Title:       "Complete Report",
			Description: "Finish the monthly financial report and submit it to management.",
			TimeStarted: time.Date(2024, 12, 07, 14, 45, 0, 0, time.UTC),
			TimeEnded:   time.Date(2024, 12, 07, 15, 15, 0, 0, time.UTC),
			Tags: []models.Tag{
				{Id: 1, Name: "work"},
				{Id: 2, Name: "report"},
				{Id: 3, Name: "urgent"},
			},
		},
		{
			ID:          "2",
			Title:       "Plan Vacation",
			Description: "Research and book flights and hotels for the summer vacation.",
			TimeStarted: time.Date(2024, 12, 07, 15, 15, 0, 0, time.UTC),
			TimeEnded:   time.Date(2024, 12, 07, 16, 15, 0, 0, time.UTC),
			Tags: []models.Tag{
				{Id: 4, Name: "personal"},
				{Id: 5, Name: "travel"},
			},
		},
		{
			ID:          "3",
			Title:       "Fix Database Issue",
			Description: "Investigate and resolve the database connectivity issue.",
			TimeStarted: time.Date(2024, 12, 07, 16, 25, 0, 0, time.UTC),
			TimeEnded:   time.Date(2024, 12, 07, 16, 45, 0, 0, time.UTC),
			Tags: []models.Tag{
				{Id: 1, Name: "work"},
				{Id: 6, Name: "bugfix"},
				{Id: 7, Name: "database"},
			},
		},
	}
	data := map[string]interface{}{
		"Title": "Tasks List",
		"Tasks": tasks,
	}

	h.BaseHandler(w, r, data)
}
func (h *Handler) TaskSaveHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	h.tmpl, err = template.ParseFiles("templates/base.html", "templates/header.html", "templates/footer.html", "templates/edit.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	description := "test"
	starttime := time.Now()
	endtime := time.Now()

	fmt.Println(id, title, description, starttime, endtime)

	var taskID int
	if id == "" {
		err = h.db.QueryRow("INSERT INTO tasks (title, description, starttime, endtime) VALUES ($1, $2, $3, $4) RETURNING id",
			title, description, starttime, endtime).Scan(&taskID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = h.db.Exec("UPDATE tasks SET title = $1 WHERE id = $2", title, id)
		if err != nil {
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
			return
		}
		taskID, _ = strconv.Atoi(id)
	}

	// Redirect to the task detail page
	redirectURL := h.getHrefByRouteName("tasks_list").String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// TaskViewHandler handles the "/task/{id}" route
func (h *Handler) TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	taskID := "" // Default value

	if idStr, exists := vars["id"]; exists {
		if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
			taskID = idStr
		}
	}

	h.tmpl, err = template.ParseFiles("templates/base.html", "templates/header.html", "templates/footer.html", "templates/edit.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//fmt.Println(h.getHrefByRouteName("task_save").String())
	data := map[string]interface{}{
		"Title":   "Task edit form",
		"TaskID":  taskID,
		"Message": "Here are the details of your task.",
		//"SaveFormUrl": h.getHrefByRouteName("task_save").String(),
		"SaveFormUrl": "/task/save",
	}
	h.BaseHandler(w, r, data)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	h.BaseHandler(w, r, data)
}

// Helper function to get URLs for routes
func (h *Handler) getHrefByRouteName(routeName string) *url.URL {
	href, err := h.router.Get(routeName).URL()
	if err != nil {
		log.Fatalf("Cannot create route: %v", routeName)
	}
	return href
}
