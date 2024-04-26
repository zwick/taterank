package data

import (
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DuplicateSlugError struct{}

func (e *DuplicateSlugError) Error() string {
	return "duplicate slug, unable to create"
}

type Models struct {
	Taters TaterModel
}

func GetModels(db *dynamodb.Client, logger *slog.Logger) Models {
	return Models{
		Taters: TaterModel{
			DB:     db,
			Logger: logger,
		},
	}
}
