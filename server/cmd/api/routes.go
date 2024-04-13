package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Taters
	router.HandlerFunc(http.MethodGet, "/api/v1/taters/:id", app.getTater)
	router.HandlerFunc(http.MethodPut, "/api/v1/taters/:id", app.updateTater)
	router.HandlerFunc(http.MethodGet, "/api/v1/taters", app.listTaters)

	// Rankings
	router.HandlerFunc(http.MethodGet, "/api/v1/rankings", app.listRankings)
	router.HandlerFunc(http.MethodPost, "/api/v1/rankings", app.createRanking)

	// Healthcheck
	router.HandlerFunc(http.MethodGet, "/api/v1/ping", app.healthCheck)

	return app.requestLogger(router)
}
