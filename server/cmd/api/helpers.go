package main

import (
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	stack := string(debug.Stack())

	app.logger.Error(err.Error(), "method", method, "uri", uri, "stack", stack)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
