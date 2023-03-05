package main

import (
	"encoding/json"
	"net/http"
)

// healthcheckHandler writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	// Append a newline to the JSON for readability in terminal.
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
