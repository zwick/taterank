package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"taterank.com/internal/database"
	"taterank.com/internal/models"
)

type application struct {
	logger *slog.Logger
	db     *dynamodb.Client
	taters *models.TaterModel
}

func main() {
	addr := flag.String("addr", ":3030", "HTTP network address")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := database.GetDynamoDBClient()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		db:     db,
		taters: &models.TaterModel{DB: db},
	}

	logger.Info("Starting server", "addr", *addr)

	server := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		logger.Error(err.Error(), "addr", *addr)
	} else {
		logger.Error("Something went terribly wrong. Stopping server", "addr", *addr)
	}

	os.Exit(1)
}
