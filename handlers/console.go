package handlers

import (
	"html/template"
	"net/http"
)

type Page struct {
	Result string
}

func ViewConsole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/console.html.tmpl")
	tmpl.Execute(w, nil)
}
