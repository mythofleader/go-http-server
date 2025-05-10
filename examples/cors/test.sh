#!/bin/bash

# Build the example
go build -o cors_example

# Start the server in the background
./cors_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the help endpoint
echo "Testing the help endpoint..."
curl -s http://localhost:8080/

# Test the API endpoint without an Origin header (should work)
echo -e "\n\nTesting the API endpoint without an Origin header (should work)..."
curl -v http://localhost:8080/api/data

# Test the API endpoint with an Origin header (should work with default config)
echo -e "\n\nTesting the API endpoint with an Origin header (should work with default config)..."
curl -v -H "Origin: http://example.com" http://localhost:8080/api/data

# Test a preflight request (OPTIONS)
echo -e "\n\nTesting a preflight request (OPTIONS)..."
curl -v -X OPTIONS -H "Origin: http://example.com" -H "Access-Control-Request-Method: GET" http://localhost:8080/api/data

# Kill the server
kill $SERVER_PID

# Now modify the example to use specific allowed domains
echo -e "\n\nModifying the example to use specific allowed domains..."
sed -i 's/\/\/ AllowedDomains: \[\]string{/AllowedDomains: []string{/' cors_example.go
sed -i 's/\/\/     "http:\/\/localhost:3000",/    "http:\/\/localhost:3000",/' cors_example.go
sed -i 's/\/\/     "https:\/\/example.com",/    "https:\/\/example.com",/' cors_example.go

# Build the modified example
go build -o cors_example_restricted

# Start the server in the background
./cors_example_restricted &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test the API endpoint with an allowed Origin header (should work)
echo -e "\n\nTesting the API endpoint with an allowed Origin header (should work)..."
curl -v -H "Origin: http://localhost:3000" http://localhost:8080/api/data

# Test the API endpoint with a non-allowed Origin header (should not include CORS headers)
echo -e "\n\nTesting the API endpoint with a non-allowed Origin header (should not include CORS headers)..."
curl -v -H "Origin: http://not-allowed.com" http://localhost:8080/api/data

# Kill the server
kill $SERVER_PID

# Clean up
rm cors_example cors_example_restricted