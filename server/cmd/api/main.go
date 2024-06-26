package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"taterank.com/internal/data"
	"taterank.com/internal/database"
)

type application struct {
	logger  *slog.Logger
	db      *dynamodb.Client
	models  data.Models
	appMode AppMode
}

func main() {
	addr := flag.String("addr", ":3030", "HTTP network address")
	appMode := flag.String("app_mode", os.Getenv("APP_MODE"), "Application mode (http or lambda)")

	flag.Parse()

	// If appMode is empty, default to "http"
	if *appMode == "" {
		*appMode = "http"
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := database.GetDynamoDBClient()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:  logger,
		db:      db,
		models:  data.GetModels(db, logger),
		appMode: AppMode(*appMode),
	}

	logger.Info("Starting server", "addr", *addr)

	switch app.appMode {
	case HTTP:
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
		}

		os.Exit(1)
	case Lambda:
		lambda.Start(httpadapter.New(app.routes()).ProxyWithContext)

	default:
		logger.Error("Unknown app mode", "appMode", app.appMode)
	}
}
