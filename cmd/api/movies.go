package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lsjoeberg/greenlight/internal/data"
)

// createMovieHandler handles the "POST /v1/movies" endpoint.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// showMovieHandler handles the "GET /v1/movies/:id" endpoint.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Check if movie ID is valid.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
