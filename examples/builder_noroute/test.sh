#!/bin/bash

# Test script for the ServerBuilder NoRoute and NoMethod example

# Build the example
go build -o builder_noroute_example

# Run the example
./builder_noroute_example

# Clean up
rm builder_noroute_example

echo "Test completed."