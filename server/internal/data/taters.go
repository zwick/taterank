package data

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sio/coolname"
)

var TableName = "Taterank-dev"
var PK = "Category#Potatoes"
var TaterPreparationsPrefix = "Preparation#"

type Tater struct {
	ID string `json:"id"`
	TaterFields
}

type TaterFields struct {
	Name        *string `json:"name" dynamodbav:",omitempty"`
	Description *string `json:"description" dynamodbav:",omitempty"`
}

type TaterModel struct {
	DB     *dynamodb.Client
	Logger *slog.Logger
}

// Returns a list of all taters
func (m *TaterModel) List() ([]*Tater, error) {
	keyExpression := expression.Key("PK").Equal(expression.Value(PK)).And(expression.Key("ID").BeginsWith(TaterPreparationsPrefix))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()

	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
	}

	result, err := m.DB.Query(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	taters := []*Tater{}

	attributevalue.UnmarshalListOfMaps(result.Items, &taters)
	collectionSanitizer(taters)

	return taters, nil
}

// Retrieves a tater by ID
func (m *TaterModel) Get(id string) (*Tater, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: PK},
			"ID": &types.AttributeValueMemberS{Value: TaterPreparationsPrefix + id},
		},
	}

	result, err := m.DB.GetItem(context.TODO(), input)

	if result == nil {
		return nil, fmt.Errorf("error getting tater")
	}

	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, nil
	}

	tater := Tater{}

	attributevalue.UnmarshalMap(result.Item, &tater)
	sanitizer(&tater)

	return &tater, nil
}

// Updates a tater by ID
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
			"ID": &types.AttributeValueMemberS{Value: TaterPreparationsPrefix + id},
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

// Creates a new tater
func (m *TaterModel) Create(fields TaterFields) (*string, error) {
	av, err := attributevalue.MarshalMap(fields)

	if err != nil {
		return nil, err
	}

	retryLimit := 3
	retries := 0

	var slug string

	for retries < retryLimit {
		retries++

		slug, err = coolname.SlugN(3)

		if err != nil {
			return nil, err
		}

		item := map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: PK},
			"ID": &types.AttributeValueMemberS{Value: TaterPreparationsPrefix + slug},
		}

		maps.Copy(item, av)

		putInput := &dynamodb.PutItemInput{
			TableName:           aws.String(TableName),
			Item:                item,
			ConditionExpression: aws.String("attribute_not_exists(ID)"),
			ReturnValues:        types.ReturnValueAllOld,
		}

		_, err = m.DB.PutItem(context.TODO(), putInput)

		if err != nil {
			var ccf *types.ConditionalCheckFailedException
			if errors.As(err, &ccf) {

				m.Logger.Warn(fmt.Sprintf("duplicate slug: '%v' retrying", slug), "retries", retries)

				if retries >= retryLimit {
					return nil, &DuplicateSlugError{}
				}

				continue
			} else {
				return nil, err
			}
		}

		break
	}

	return &slug, nil
}

// Given a Tater, this function will remove the TaterPreparationsPrefix from the ID
func sanitizer(data *Tater) {
	data.ID = strings.Split(data.ID, TaterPreparationsPrefix)[1]
}

// Given a collection of Taters, this function will remove the TaterPreparationsPrefix from the ID
func collectionSanitizer(data []*Tater) {
	for _, tater := range data {
		sanitizer(tater)
	}
}
