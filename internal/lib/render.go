package lib

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// RenderTemplate renders a single HTML file
func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	// Try multiple possible template paths
	templatePaths := []string{
		filepath.Join("../../web/templates", page+".html"), // Original relative path
		filepath.Join("web/templates", page+".html"),       // From root directory
		filepath.Join("./web/templates", page+".html"),     // Another possibility
	}

	var (
		tmpl          *template.Template
		err           error
		templateFound bool = false
	)

	// Try each path until we find one that works
	for _, path := range templatePaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("[INFO] Using template path: %s", path)
			tmpl, err = template.ParseFiles(path)
			if err == nil {
				templateFound = true
				break
			}
		}
	}

	// If no template found, return error
	if !templateFound {
		log.Printf("[ERROR] Failed to load template '%s'. Tried paths: %v", page, templatePaths)
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
