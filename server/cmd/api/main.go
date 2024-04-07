package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"taterank.com/internal/models"
)

type application struct {
	logger *slog.Logger
	db     *dynamodb.Client
	taters *models.TaterModel
}

func getDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	return client, nil
}

func main() {
	addr := flag.String("addr", ":3030", "HTTP network address")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := getDynamoDBClient()

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

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error(), "addr", *addr)
	os.Exit(1)
}
