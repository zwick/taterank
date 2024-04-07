#!/bin/bash

# Set the AWS CLI Docker image
AWS_CLI_IMAGE="amazon/aws-cli"

# Set the AWS region and credentials
AWS_REGION="us-east-1"
AWS_ACCESS_KEY_ID="test"
AWS_SECRET_ACCESS_KEY="test"

# Set the custom endpoint URL
CUSTOM_ENDPOINT="http://localhost:4566"

# Set the data file path
DATA_FILE="$(pwd)/seeds/dynamodb-seed-data.json"
CREDS_FILE="$(pwd)/development/credentials"

# Generate seed filed

# Run the AWS CLI Docker container and execute the batch write command
docker run --rm \
  --net=host \
  -v $DATA_FILE:/data.json \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  $AWS_CLI_IMAGE dynamodb batch-write-item --request-items file:///data.json \
  --endpoint-url $CUSTOM_ENDPOINT \
  --region=$AWS_REGION

echo "Seeding dynamodb complete!"