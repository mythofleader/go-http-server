#!/bin/bash

# Test script for the middleware flow control example

# Build and run the example
go build -o middleware-example
./middleware-example &
PID=$!

# Wait for the server to start
sleep 2

# Make requests to different endpoints
echo "Testing / endpoint..."
curl -s http://localhost:8080/
echo -e "\n"

echo "Testing /skip endpoint..."
curl -s http://localhost:8080/skip
echo -e "\n"

echo "Testing /order endpoint..."
curl -s http://localhost:8080/order
echo -e "\n"

echo "Testing /help endpoint..."
curl -s http://localhost:8080/help
echo -e "\n"

# Kill the server
kill $PID

echo "Test completed. Check the console output to see the logs."