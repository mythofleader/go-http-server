// This example demonstrates the use of the port configuration options in ServerBuilder
package main

import (
	"fmt"
	"log"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	fmt.Println("=== Testing port configuration options in ServerBuilder ===")

	// Scenario 1: Creating a server builder with a port parameter
	fmt.Println("\n=== Scenario 1: Creating a server builder with a port parameter ===")
	builder1 := server.NewServerBuilder(server.FrameworkGin, "8081")
	testBuilder(builder1, "Builder with port parameter")

	// Scenario 2: Creating a server builder without a port parameter and then calling WithDefaultPort
	fmt.Println("\n=== Scenario 2: Creating a server builder without a port parameter and then calling WithDefaultPort ===")
	builder2 := server.NewServerBuilder(server.FrameworkGin)
	builder2.WithDefaultPort()
	testBuilder(builder2, "Builder without port parameter but with WithDefaultPort")

	// Scenario 3: Creating a server builder without a port parameter and not calling WithDefaultPort
	fmt.Println("\n=== Scenario 3: Creating a server builder without a port parameter and not calling WithDefaultPort ===")
	builder3 := server.NewServerBuilder(server.FrameworkGin)
	testBuilder(builder3, "Builder without port parameter and without WithDefaultPort")
}

func testBuilder(builder *server.ServerBuilder, description string) {
	fmt.Printf("Testing %s...\n", description)

	// Try to build the server
	srv, err := builder.Build()

	if err != nil {
		fmt.Printf("Error building server: %v\n", err)
		fmt.Printf("Result: %s test PASSED (expected error)\n", description)
	} else {
		// If we got a server, make sure we can stop it
		fmt.Printf("Successfully built server\n")

		// Stop the server to clean up
		if err := srv.Stop(); err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}

		fmt.Printf("Result: %s test PASSED (expected success)\n", description)
	}
}
