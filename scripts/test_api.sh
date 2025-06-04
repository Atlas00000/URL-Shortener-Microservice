#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing URL Shortener API..."

# Test health endpoint
echo -e "\n${GREEN}Testing health endpoint...${NC}"
curl -s http://localhost:8080/health | jq '.'

# Test URL shortening
echo -e "\n${GREEN}Testing URL shortening...${NC}"
SHORT_URL_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.google.com","expiration_days":30}' \
  http://localhost:8080/shorten)

echo $SHORT_URL_RESPONSE | jq '.'

# Extract short ID from response
SHORT_ID=$(echo $SHORT_URL_RESPONSE | jq -r '.short_url')

if [ "$SHORT_ID" != "null" ]; then
  # Test redirection
  echo -e "\n${GREEN}Testing redirection...${NC}"
  curl -s -I http://localhost:8080/$SHORT_ID

  # Test click recording
  echo -e "\n${GREEN}Testing click recording...${NC}"
  curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{\"short_id\":\"$SHORT_ID\",\"ip\":\"127.0.0.1\",\"user_agent\":\"curl/7.68.0\"}" \
    http://localhost:8080/analytics/click | jq '.'

  # Test analytics
  echo -e "\n${GREEN}Testing analytics...${NC}"
  curl -s "http://localhost:8080/analytics?short_id=$SHORT_ID" | jq '.'
else
  echo -e "\n${RED}Failed to get short URL. Skipping subsequent tests.${NC}"
fi 