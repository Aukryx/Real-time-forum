package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"lib"
	"log"
	"models"
	"net"
	"net/http"
	"runtime/debug"
)

// HandleError renders an error page or sends JSON response if API request
func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" {
		http.Error(w, message, statusCode)
		return
	}

	data := models.PageData{
		Title:     "Error",
		Header:    fmt.Sprintf("Error %d", statusCode),
		Content:   map[string]template.HTML{"Msg_raw": template.HTML("<h1>" + message + "</h1>")},
		IsError:   true,
		ErrorCode: statusCode,
	}

	lib.RenderTemplate(w, "index", data)
}

// WithErrorHandling middleware catches panics and handles errors
func WithErrorHandling(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERROR] %v\n%s", err, debug.Stack())

				// Default to 500 Internal Server Error
				statusCode := http.StatusInternalServerError
				message := "Internal Server Error"

				// Check error type
				switch e := err.(type) {
				case *models.CustomError:
					statusCode = e.StatusCode
					message = e.Message
				case *net.OpError:
					statusCode = http.StatusInternalServerError
					message = "A network error occurred"
				case string:
					// Handle direct string panics (should be avoided)
					switch e {
					case "bad request":
						statusCode = http.StatusBadRequest
						message = "Bad Request"
					case "not found":
						statusCode = http.StatusNotFound
						message = "Not Found"
					default:
						statusCode = http.StatusInternalServerError
						message = "An unexpected error occurred"
					}
				default:
					// Attempt to match error using errors.As()
					var customErr *models.CustomError
					if errors.As(err.(error), &customErr) {
						statusCode = customErr.StatusCode
						message = customErr.Message
					}
				}

				// Send notification (optional)
				// lib.PostItOnNfty("[FORUM server]", message)

				// Render error page or JSON response
				HandleError(w, r, statusCode, message)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
