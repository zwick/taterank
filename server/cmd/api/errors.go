package main

import (
	"net/http"
	"runtime/debug"
)

func (app *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	stack := string(debug.Stack())

	app.logger.Error(err.Error(), "method", method, "uri", uri, "stack", stack)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, payload{"error": message}, status, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, payload{"error": message})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

