package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Home handles the home page request
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Get the absolute path to the HTML file
	homePath := filepath.Join("web", "templates", "index.html")

	// Parse the HTML file
	tmpl, err := template.ParseFiles(homePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with any data needed
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
