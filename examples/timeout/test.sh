#!/bin/bash

# Test script for the timeout middleware example

# Build and run the example
go build -o timeout-example
./timeout-example &
PID=$!

# Wait for the server to start
sleep 2

# Make requests to different endpoints
echo "Testing / endpoint (immediate response)..."
curl -s http://localhost:8080/
echo -e "\n"

echo "Testing /delay/1 endpoint (1-second delay)..."
curl -s http://localhost:8080/delay/1
echo -e "\n"

echo "Testing /delay/2 endpoint (2-second delay)..."
curl -s http://localhost:8080/delay/2
echo -e "\n"

echo "Testing /timeout endpoint (should timeout after 3 seconds)..."
time curl -s http://localhost:8080/timeout
echo -e "\n"

echo "Testing /help endpoint..."
curl -s http://localhost:8080/help
echo -e "\n"

# Kill the server
kill $PID

echo "Test completed. The /timeout endpoint should have returned a 503 Service Unavailable response after 3 seconds."