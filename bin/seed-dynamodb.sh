#!/bin/bash

# Set the AWS CLI Docker image
AWS_CLI_IMAGE="amazon/aws-cli"

# Set the AWS region and credentials
AWS_REGION="us-east-1"
AWS_ACCESS_KEY_ID="test"
AWS_SECRET_ACCESS_KEY="test"

# Set the DynamoDB table name
TABLE_NAME="Taterank-dev"
# Set the custom endpoint URL
CUSTOM_ENDPOINT="http://localhost:4566"

# Set the data file path
DATA_FILE="$(pwd)/seeds/dynamodb-seed-data.json"

# Generate seed file


# Run the AWS CLI Docker container and execute the batch write command
docker run --rm -v $DATA_FILE:/data.json $AWS_CLI_IMAGE dynamodb batch-write-item \
  --region $AWS_REGION \
  --endpoint-url $CUSTOM_ENDPOINT \
  --request-items file:///data.json
