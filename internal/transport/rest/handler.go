package rest

import (
	"crud_test/internal/models"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
	tmpl   *template.Template
}

func NewHandler(router *mux.Router) *Handler {
	// funcMap := template.FuncMap{
	// 	"formatTime": formatTimestamp, // Registering the function
	// }

	// tmpl := template.New("").Funcs(funcMap) // Initialize with FuncMap
	tmpl := template.New("")
	return &Handler{
		router: router,
		tmpl:   tmpl,
	}
}

// BaseHandler renders the base layout and dynamically includes specific content
func (h *Handler) BaseHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {

	w.WriteHeader(http.StatusOK)

	data["Home_page_url"] = h.getHrefByRouteName("home")
	data["Tasks_page_url"] = h.getHrefByRouteName("tasks_list")
	data["Sign_page_url"] = h.getHrefByRouteName("login")

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

func (h *Handler) TaskViewHandlerzz(w http.ResponseWriter, r *http.Request) {
	// dummy
}

// TasksListHandler handles the "/tasks" route
func (h *Handler) TasksListHandler(w http.ResponseWriter, r *http.Request) {
	// Register the FuncMap before parsing templates
	var err error
	funcMap := template.FuncMap{
		"formatTime": formatTimestamp, // Registering the function
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
			TimeStarted: 1702000000,
			TimeEnded:   1702010000,
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
			TimeStarted: 1701990000,
			TimeEnded:   1701995000,
			Tags: []models.Tag{
				{Id: 4, Name: "personal"},
				{Id: 5, Name: "travel"},
			},
		},
		{
			ID:          "3",
			Title:       "Fix Database Issue",
			Description: "Investigate and resolve the database connectivity issue.",
			TimeStarted: 1701980000,
			TimeEnded:   1701988000,
			Tags: []models.Tag{
				{Id: 1, Name: "work"},
				{Id: 6, Name: "bugfix"},
				{Id: 7, Name: "database"},
			},
		},
	}
	data := map[string]interface{}{
		"Title":       "Tasks List",
		"Tasks":       tasks,
		"taskViewUrl": h.getHrefByRouteName("task_view_zz"),
	}

	h.BaseHandler(w, r, data)
}

// TaskViewHandler handles the "/task/{id}" route
func (h *Handler) TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r) // Get URL variables
	taskID := vars["id"]

	h.tmpl, err = template.ParseFiles("templates/base.html", "templates/header.html", "templates/footer.html", "templates/edit.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"Title":   "Task Details",
		"TaskID":  taskID,
		"Message": "Here are the details of your task.",
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
