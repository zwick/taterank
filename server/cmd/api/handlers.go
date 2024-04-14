package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"taterank.com/internal/models"
)

// Get a tater by ID
func (app *application) getTater(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	taterID := params.ByName("id")

	if taterID == "" {
		http.NotFound(w, r)
		return
	}

	tater, err := app.taters.Get(taterID)

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
	taters, err := app.taters.List()

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
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	if id == "" {
		http.NotFound(w, r)
		return
	}

	var fields models.TaterFields

	err := app.readJSON(r, &fields)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.taters.Update(id, fields)

	if err != nil {
		app.logError(r, err)
		app.errorResponse(w, r, http.StatusBadRequest, "failed to update tater")
		return
	}

	updatedTater := models.Tater{
		ID:          id,
		TaterFields: fields,
	}

	app.writeJSON(w, payload{"data": updatedTater}, http.StatusOK, nil)
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
