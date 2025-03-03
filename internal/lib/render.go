package lib

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Base directory for HTML files
const TemplateDir = "../../web/templates"

// RenderTemplate renders a single HTML file
func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	// Define the template file path
	pagePath := filepath.Join(TemplateDir, page+".html")

	// Parse the template
	tmpl, err := template.ParseFiles(pagePath)
	if err != nil {
		log.Printf("[ERROR] Failed to load template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with data
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("[ERROR] Template execution failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
