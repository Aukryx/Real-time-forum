package handlers

import (
	"lib"
	"models"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:   "Home",
		Header:  "Welcome to the Forum",
		Content: "This is the homepage.",
	}
	lib.RenderTemplate(w, "index", data)
}
