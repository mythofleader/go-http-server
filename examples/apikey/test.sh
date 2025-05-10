#!/bin/bash

# Build the example
go build -o apikey_example

# Start the server in the background
./apikey_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the help endpoint
echo "Testing the help endpoint..."
curl -s http://localhost:8080/help

# Test the protected endpoint without an API key (should fail)
echo -e "\n\nTesting the protected endpoint without an API key (should fail)..."
curl -v http://localhost:8080/api/data

# Test the protected endpoint with an invalid API key (should fail)
echo -e "\n\nTesting the protected endpoint with an invalid API key (should fail)..."
curl -v -H "x-api-key: wrong-api-key" http://localhost:8080/api/data

# Test the protected endpoint with a valid API key (should succeed)
echo -e "\n\nTesting the protected endpoint with a valid API key (should succeed)..."
curl -v -H "x-api-key: my-secret-api-key" http://localhost:8080/api/data

# Kill the server
kill $SERVER_PID

# Clean up
rm apikey_example