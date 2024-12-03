package rest

import (
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
	tmpl   *template.Template
}

func NewHandler(router *mux.Router) *Handler {
	// Parse all templates in the "templates/" directory
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &Handler{
		router: router,
		tmpl:   tmpl,
	}
}

// BaseHandler renders the base layout and dynamically includes specific content
func (h *Handler) BaseHandler(w http.ResponseWriter, r *http.Request, content string, data map[string]interface{}) {
	w.WriteHeader(http.StatusOK)

	// Dynamically include content template name
	data["Content"] = content

	// Render the base layout (e.g., index.html)
	err := h.tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HomeHandler handles the "/" route and injects specific content
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "Welcome to simple task tracker",
		"Description": "Manage your tasks efficiently",
		"newTaskHref": h.getHrefByRouteName("task_new").String(),
	}
	err := h.tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	h.BaseHandler(w, r, "home.html", data)
}

// TasksListHandler handles the "/tasks" route
func (h *Handler) TasksListHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Tasks List",
		"Tasks": []string{"Task 1", "Task 2", "Task 3"},
	}
	h.BaseHandler(w, r, "tasks_list.html", data)
}

// TaskViewHandler handles the "/task/{id}" route
func (h *Handler) TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get URL variables
	taskID := vars["id"]

	data := map[string]interface{}{
		"Title":   "Task Details",
		"TaskID":  taskID,
		"Message": "Here are the details of your task.",
	}
	h.BaseHandler(w, r, "task/edit.html", data)
}

// Helper function to get URLs for routes
func (h *Handler) getHrefByRouteName(routeName string) *url.URL {
	href, err := h.router.Get(routeName).URL()
	if err != nil {
		log.Fatalf("Cannot create route: %v", routeName)
	}
	return href
}
