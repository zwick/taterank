package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Returns a new DynamoDB client using the default configuration
// and credentials provider chain.
func GetDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	_, err = client.Options().Credentials.Retrieve(context.TODO())

	if err != nil {
		return nil, err
	}

	return client, nil
}

type TestConfigOptions struct {
	Endpoint string
}

// Returns a new DynamoDB client for testing.
func GetTestDynamoDBClient(t *testing.T, options TestConfigOptions) (*dynamodb.Client, error) {
	endpoint := "http://localhost:4566"

	if options.Endpoint != "" {
		endpoint = options.Endpoint
	}

	region := "us-east-1"

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		t.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg)

	_, err = client.Options().Credentials.Retrieve(context.TODO())

	if err != nil {
		t.Fatal(err)
	}

	return client, nil
}
