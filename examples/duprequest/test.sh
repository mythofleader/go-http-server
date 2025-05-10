#!/bin/bash

# Build the example
go build -o duprequest_example

# Start the server in the background
./duprequest_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the help endpoint
echo "Testing the help endpoint..."
curl -s http://localhost:8080/

# Test creating an order (should succeed)
echo -e "\n\nTesting order creation (should succeed)..."
curl -v -X POST -d '{"id":"123","product":"example"}' http://localhost:8080/api/orders

# Test creating the same order again (should fail with 409 Conflict)
echo -e "\n\nTesting duplicate order creation (should fail with 409 Conflict)..."
curl -v -X POST -d '{"id":"123","product":"example"}' http://localhost:8080/api/orders

# Test creating a different order (should succeed)
echo -e "\n\nTesting different order creation (should succeed)..."
curl -v -X POST -d '{"id":"124","product":"example"}' http://localhost:8080/api/orders

# Kill the server
kill $SERVER_PID

# Clean up
rm duprequest_example

echo -e "\n\nTest completed."