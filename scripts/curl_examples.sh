#!/bin/bash

# Configuration
DOMAIN="https://api.simpleaiwork.com"
# DOMAIN="http://localhost:8080" # Use this for local testing

echo "==========================================="
echo "Target: $DOMAIN"
echo "==========================================="

# 1. Insert User (POST)
# We only need to send 'email'; 'email_norm' is handled by the server.
echo -e "\n[POST] Creating user..."
curl -s -X POST "$DOMAIN/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "DemoUser@SimpleAIWork.com",
    "status": 1
  }' | jq . || echo " (install jq for pretty output)"

# 2. Query User (GET)
# The server will normalize the email query param.
echo -e "\n\n[GET] Querying user..."
curl -s "$DOMAIN/users?email=DemoUser@SimpleAIWork.com" | jq . || echo " (install jq for pretty output)"

echo -e "\n\nDone."
