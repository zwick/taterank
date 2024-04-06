package main

import (
	"log"
	"net/http"
)

// TODO: Use gql duh!

func listRankings(w http.ResponseWriter, r *http.Request) {
	// Get top rankings
	w.Write([]byte("Hello, World!"))
}

func listRecentRankings(w http.ResponseWriter, r *http.Request) {

}

func updateRankingVotes(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/rankings/{$}", listRankings)
	mux.HandleFunc("GET /", listRankings)

	log.Print("Starting server on :3030")

	err := http.ListenAndServe(":3030", mux)
	log.Fatal(err)
}
