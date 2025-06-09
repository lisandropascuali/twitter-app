#!/bin/bash

# Set AWS credentials for LocalStack
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1

# Create DynamoDB table for tweets
aws dynamodb create-table \
    --endpoint-url http://localhost:4566 \
    --region us-east-1 \
    --table-name tweets \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=user_id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --global-secondary-indexes \
        "[
            {
                \"IndexName\": \"user_id-index\",
                \"KeySchema\": [
                    {\"AttributeName\":\"user_id\",\"KeyType\":\"HASH\"}
                ],
                \"Projection\": {
                    \"ProjectionType\":\"ALL\"
                },
                \"ProvisionedThroughput\": {
                    \"ReadCapacityUnits\": 5,
                    \"WriteCapacityUnits\": 5
                }
            }
        ]" \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5

# Verify table creation
echo "Verifying table creation..."
aws dynamodb describe-table \
    --endpoint-url http://localhost:4566 \
    --region us-east-1 \
    --table-name tweets 