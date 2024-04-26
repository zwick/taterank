package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"taterank.com/internal/data"
)

// Get a tater by ID
func (app *application) getTater(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	taterID := params.ByName("id")

	if taterID == "" {
		app.notFoundResponse(w, r)
		return
	}

	tater, err := app.models.Taters.Get(taterID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if tater == nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, payload{"data": tater}, http.StatusOK, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Create tater
func (app *application) createTater(w http.ResponseWriter, r *http.Request) {

	var fields data.TaterFields

	err := app.readJSON(r, &fields)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	slug, err := app.models.Taters.Create(fields)

	if err != nil {
		var ds *data.DuplicateSlugError
		if errors.As(err, &ds) {
			app.serverErrorResponse(w, r, err)
			return
		}

		app.logError(r, err)
		app.errorResponse(w, r, http.StatusBadRequest, "failed to update")
		return
	}

	response := struct {
		ID *string `json:"id"`
	}{
		ID: slug,
	}

	app.writeJSON(w, payload{"data": response}, http.StatusOK, nil)
}

// List all taters
func (app *application) listTaters(w http.ResponseWriter, r *http.Request) {
	taters, err := app.models.Taters.List()

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
		app.notFoundResponse(w, r)
		return
	}

	var fields data.TaterFields

	err := app.readJSON(r, &fields)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Taters.Update(id, fields)

	// TODO: handle not found error
	if err != nil {
		app.logError(r, err)
		app.errorResponse(w, r, http.StatusBadRequest, "failed to update tater")
		return
	}

	updatedTater := data.Tater{
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
