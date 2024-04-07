package main

import (
	"encoding/json"
	"net/http"
)

// Get a tater by ID
func (app *application) getTater(w http.ResponseWriter, r *http.Request) {
	tater, err := app.taters.GetByID(r.PathValue("id"))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if tater == nil {
		http.NotFound(w, r)
		return
	}

	response, err := json.Marshal(tater)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (app *application) listRankings(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte{})
}

func (app *application) createRanking(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte{})
}
