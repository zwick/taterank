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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaterModel struct {
	DB *dynamodb.Client
}

func (m *TaterModel) Get() ([]*Tater, error) {
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

func (m *TaterModel) Update(tater Tater) (*Tater, error) {
	update := expression.Set(expression.Name("Description"), expression.Value(tater.Description))
	update.Set(expression.Name("Name"), expression.Value(tater.Name))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: PK},
			"ID": &types.AttributeValueMemberS{Value: tater.ID},
		},
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	}

	output, err := m.DB.UpdateItem(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	updatedTater := &Tater{}
	attributevalue.UnmarshalMap(output.Attributes, &updatedTater)

	return updatedTater, nil
}
