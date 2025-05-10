#!/bin/bash

# Build the example
go build -o errors_example

# Start the server in the background and redirect output to a log file
./errors_example > server.log 2>&1 &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the help endpoint
echo "Testing the help endpoint..."
curl -s http://localhost:8080/

# Test each error endpoint
echo -e "\n\nTesting the bad-request endpoint..."
curl -v http://localhost:8080/bad-request

echo -e "\n\nTesting the unauthorized endpoint..."
curl -v http://localhost:8080/unauthorized

echo -e "\n\nTesting the forbidden endpoint..."
curl -v http://localhost:8080/forbidden

echo -e "\n\nTesting the not-found endpoint..."
curl -v http://localhost:8080/not-found

echo -e "\n\nTesting the conflict endpoint..."
curl -v http://localhost:8080/conflict

echo -e "\n\nTesting the internal-error endpoint..."
curl -v http://localhost:8080/internal-error

echo -e "\n\nTesting the service-unavailable endpoint..."
curl -v http://localhost:8080/service-unavailable

echo -e "\n\nTesting the custom-error endpoint..."
curl -v http://localhost:8080/custom-error

echo -e "\n\nTesting the from-http-error endpoint..."
curl -v http://localhost:8080/from-http-error

echo -e "\n\nTesting the error-method endpoint..."
curl -v http://localhost:8080/error-method

echo -e "\n\nTesting the multiple-errors endpoint..."
curl -v http://localhost:8080/multiple-errors

echo -e "\n\nTesting the invalid-request-param endpoint..."
curl -v http://localhost:8080/invalid-request-param

# Kill the server
kill $SERVER_PID

# Clean up
rm errors_example

echo -e "\n\nTest completed."
