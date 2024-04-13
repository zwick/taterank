package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Get a tater by ID
func (app *application) getTater(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	taterID := params.ByName("id")

	if taterID == "" {
		http.NotFound(w, r)
		return
	}

	tater, err := app.taters.GetByID(taterID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if tater == nil {
		http.NotFound(w, r)
		return
	}

	err = app.writeJSON(w, payload{"data": tater}, http.StatusOK, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// List all taters
func (app *application) listTaters(w http.ResponseWriter, r *http.Request) {
	taters, err := app.taters.Get()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, payload{"data": taters}, http.StatusOK, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Update tater
func (app *application) updateTater(w http.ResponseWriter, r *http.Request) {
	// params := httprouter.ParamsFromContext(r.Context())

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := app.readJSON(r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.writeJSON(w, payload{"data": input}, http.StatusOK, nil)
}

func (app *application) listRankings(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte{})
}

func (app *application) createRanking(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte{})
}

func (app *application) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}
