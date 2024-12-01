package rest

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var tmpl *template.Template

func init() {
	// Parse all templates in the "templates/" directory
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func HomeHandler(w http.ResponseWriter, r *http.Request, router *mux.Router) {
	w.WriteHeader(http.StatusOK)

	//newTaskHref := r.Get("/task/new").URL("category", "technology", "id", "42")
	newTaskHref, err := router.Get("task_new").URL()
	if err != nil {
		log.Fatal("Can`t create new task route")
	}
	data := map[string]string{
		"Title":       "Welcome to simple task tracker",
		"Description": "Bla Bla Bla",
		"newTaskHref": newTaskHref.String(),
	}
	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TasksListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Tasks list here")
}

func TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task selected id: %v\n", vars["id"])
}
