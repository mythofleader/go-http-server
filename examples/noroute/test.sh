#!/bin/bash

# Test script for the NoRoute and NoMethod example

# Build the example
go build -o noroute-example

# Start the server in the background
./noroute-example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the valid route
echo "Testing the valid route (GET /)..."
curl -s http://localhost:8080/
echo -e "\n"

# Test the NoRoute handler (404 Not Found)
echo "Testing the NoRoute handler (GET /nonexistent)..."
curl -v http://localhost:8080/nonexistent
echo -e "\n"

# Test the NoMethod handler (405 Method Not Allowed)
echo "Testing the NoMethod handler (POST /)..."
curl -v -X POST http://localhost:8080/
echo -e "\n"

# Kill the server
kill $SERVER_PID

# Clean up
rm noroute-example

echo "Test completed."