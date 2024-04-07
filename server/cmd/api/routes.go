package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Taters
	mux.HandleFunc("GET /api/taters/{id}", app.getTater)

	// Rankings
	mux.HandleFunc("GET /api/rankings", app.listRankings)
	mux.HandleFunc("POST /api/rankings", app.createRanking)

	return app.handlePanic(app.requestLogger(mux))
}
