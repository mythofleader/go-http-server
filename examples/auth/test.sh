#!/bin/bash

# Build the example
go build -o auth_example

# Start the server in the background
./auth_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test with JWT authentication (should succeed)
echo "Testing with JWT authentication (should succeed)..."
curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSIsIm5hbWUiOiJVc2VyIE9uZSIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjo0NjY5MzM0NDAwfQ.Kk-VYwgAsk4XVnIp7MiCHs9q6NcHYUQQwrBVW0XM7bQ" http://localhost:8080/api/profile

# Test with Basic authentication (should fail because we're using JWT authentication)
echo -e "\n\nTesting with Basic authentication (should fail because we're using JWT authentication)..."
curl -v -H "Authorization: Basic $(echo -n 'user1:password' | base64)" http://localhost:8080/api/profile

# Kill the server
kill $SERVER_PID

# Modify the example to use Basic authentication
sed -i 's/AuthType:   server.AuthTypeJWT/AuthType:   server.AuthTypeBasic/' auth_example.go
sed -i 's/\/\/ basicAuthConfig/basicAuthConfig/' auth_example.go
sed -i 's/\/\/     UserLookup/    UserLookup/' auth_example.go
sed -i 's/\/\/     AuthType/    AuthType/' auth_example.go
sed -i 's/server.AuthMiddleware(authConfig)/server.AuthMiddleware(basicAuthConfig)/' auth_example.go

# Build the modified example
go build -o auth_example_basic

# Start the server in the background
./auth_example_basic &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Test with Basic authentication (should succeed)
echo -e "\n\nTesting with Basic authentication (should succeed)..."
curl -v -H "Authorization: Basic $(echo -n 'user1:password' | base64)" http://localhost:8080/api/profile

# Test with JWT authentication (should fail because we're using Basic authentication)
echo -e "\n\nTesting with JWT authentication (should fail because we're using Basic authentication)..."
curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSIsIm5hbWUiOiJVc2VyIE9uZSIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjo0NjY5MzM0NDAwfQ.Kk-VYwgAsk4XVnIp7MiCHs9q6NcHYUQQwrBVW0XM7bQ" http://localhost:8080/api/profile

# Kill the server
kill $SERVER_PID

# Clean up
rm auth_example auth_example_basic