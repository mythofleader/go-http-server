#!/bin/bash

# Test script for the custom NoRoute and NoMethod handlers example

# Build the example
go build -o noroute_custom_example

# Start the server in the background
./noroute_custom_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the valid route
echo "Testing the valid route (GET /)..."
curl -s http://localhost:8080/
echo -e "\n"

# Test the custom NoRoute handler (404 Not Found)
echo "Testing the custom NoRoute handler (GET /nonexistent)..."
curl -v http://localhost:8080/nonexistent
echo -e "\n"

# Test the custom NoMethod handler (405 Method Not Allowed)
echo "Testing the custom NoMethod handler (POST /)..."
curl -v -X POST http://localhost:8080/
echo -e "\n"

# Kill the server
kill $SERVER_PID

# Clean up
rm noroute_custom_example

echo "Test completed."