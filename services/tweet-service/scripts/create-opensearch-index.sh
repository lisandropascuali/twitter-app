#!/bin/bash

# Wait for OpenSearch to be ready
echo "Waiting for OpenSearch to be ready..."
until curl -s http://localhost:9200/_cluster/health | grep -q '"status":"green"\|"status":"yellow"'; do
    sleep 1
done

# Create tweets index with mapping
echo "Creating tweets index..."
curl -X PUT "http://localhost:9200/tweets" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "user_id": { "type": "keyword" },
      "content": { "type": "text" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

echo "OpenSearch index created successfully!" 