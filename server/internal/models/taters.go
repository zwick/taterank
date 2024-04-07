package models

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "Taterank-dev"
var PK = "Category#Potatoes"

type Tater struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type TaterModel struct {
	DB *dynamodb.Client
}

// Retrieves a tater by ID
func (m *TaterModel) GetByID(id string) (*Tater, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: PK},
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := m.DB.GetItem(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, nil
	}

	tater := Tater{}

	attributevalue.UnmarshalMap(result.Item, &tater)

	return &tater, nil
}
