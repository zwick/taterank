package main

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Taters
	router.HandlerFunc(http.MethodGet, "/api/v1/taters/:id", app.getTater)
	router.HandlerFunc(http.MethodPut, "/api/v1/taters/:id", app.updateTater)
	router.HandlerFunc(http.MethodGet, "/api/v1/taters", app.listTaters)
	router.HandlerFunc(http.MethodPost, "/api/v1/taters", app.createTater)

	// Rankings
	router.HandlerFunc(http.MethodGet, "/api/v1/rankings", app.listRankings)
	router.HandlerFunc(http.MethodPost, "/api/v1/rankings", app.createRanking)

	// Healthcheck
	router.HandlerFunc(http.MethodGet, "/api/v1/ping", app.healthCheck)

	return gziphandler.GzipHandler(app.requestLogger(router))
}
