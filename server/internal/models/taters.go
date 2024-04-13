package models

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "Taterank-dev"
var PK = "Category#Potatoes"

type Tater struct {
	ID string `json:"id"`
	TaterFields
}

type TaterFields struct {
	Name        *string `json:"name" dynamodbav:",omitempty"`
	Description *string `json:"description" dynamodbav:",omitempty"`
}

type TaterModel struct {
	DB *dynamodb.Client
}

func (m *TaterModel) List() ([]*Tater, error) {
	keyExpression := expression.Key("PK").Equal(expression.Value(PK))
	expression, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()

	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		KeyConditionExpression:    expression.KeyCondition(),
		ExpressionAttributeValues: expression.Values(),
		ExpressionAttributeNames:  expression.Names(),
	}

	result, err := m.DB.Query(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	taters := []*Tater{}

	attributevalue.UnmarshalListOfMaps(result.Items, &taters)

	return taters, nil
}

// Retrieves a tater by ID
func (m *TaterModel) Get(id string) (*Tater, error) {
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

func (m *TaterModel) Update(id string, fields TaterFields) error {
	av, err := attributevalue.MarshalMap(fields)

	if err != nil {
		return nil
	}

	update := expression.UpdateBuilder{}

	// Build the expression dynamically using the output of MarshalMap()
	// to allow partial updates.
	for k, v := range av {
		update = update.Set(expression.Name(k), expression.Value(v))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}

	updateInput := &types.Update{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: PK},
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	}

	writeInput := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Update: updateInput,
			},
		},
	}

	_, err = m.DB.TransactWriteItems(context.TODO(), writeInput)

	if err != nil {
		return err
	}

	return nil
}
