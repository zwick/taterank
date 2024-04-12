package main

import (
	"flag"
	"github.com/aws/aws-lambda-go/lambda"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"taterank.com/internal/database"
	"taterank.com/internal/models"
)

type application struct {
	logger  *slog.Logger
	db      *dynamodb.Client
	taters  *models.TaterModel
	appMode AppMode
}

func main() {
	addr := flag.String("addr", ":3030", "HTTP network address")
	appMode := flag.String("app_mode", "lambda", "App mode (http or lambda)")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := database.GetDynamoDBClient()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:  logger,
		db:      db,
		taters:  &models.TaterModel{DB: db},
		appMode: AppMode(*appMode),
	}

	logger.Info("Starting server", "addr", *addr)

	switch app.appMode {
	case HTTP:
		server := &http.Server{
			Addr:    *addr,
			Handler: app.routes(),
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
