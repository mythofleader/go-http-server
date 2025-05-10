#!/bin/bash

# Test script for the logging middleware example

# Build and run the example
go build -o logging-example
./logging-example &
PID=$!

# Wait for the server to start
sleep 2

# Make requests to different endpoints
echo "Testing / endpoint..."
curl -s http://localhost:8081/
echo -e "\n"

echo "Testing /json endpoint..."
curl -s http://localhost:8081/json
echo -e "\n"

echo "Testing /error endpoint..."
curl -s http://localhost:8081/error
echo -e "\n"

echo "Testing /health endpoint (should not be logged)..."
curl -s http://localhost:8081/health
echo -e "\n"

echo "Testing /help endpoint..."
curl -s http://localhost:8081/help
echo -e "\n"

# Kill the server
kill $PID

echo "Test completed. Check the console output to see the logs."
