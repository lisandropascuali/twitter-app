#!/bin/bash

# Set AWS credentials for LocalStack
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1

# List all items from the tweets table
echo "Listing all tweets from DynamoDB..."
aws dynamodb scan \
    --endpoint-url http://localhost:4566 \
    --region us-east-1 \
    --table-name tweets \
    --select ALL_ATTRIBUTES

# Optional: Format the output as JSON
# aws dynamodb scan \
#     --endpoint-url http://localhost:4566 \
#     --region us-east-1 \
#     --table-name tweets \
#     --select ALL_ATTRIBUTES \
#     --output json 