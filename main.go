package main

import (
	"fmt"
	"os"
)

// main serves as a simple entry point that directs users to the correct application entry points.
// The actual applications are located in cmd/api and cmd/admin directories.
func main() {
	fmt.Println("Kratos-Kit Project")
	fmt.Println("==================")
	fmt.Println()
	fmt.Println("Available applications:")
	fmt.Println("  API Server:   go run cmd/api/main.go")
	fmt.Println("  Admin Server: go run cmd/admin/main.go")
	fmt.Println("  CLI Tool:     go run cmd/cli/main.go")
	fmt.Println()
	fmt.Println("For development with hot reload:")
	fmt.Println("  make run")
	fmt.Println()
	fmt.Println("For more commands:")
	fmt.Println("  make help")

	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println()
		fmt.Println("Build and run instructions:")
		fmt.Println("  make build    # Build all applications")
		fmt.Println("  make test     # Run tests")
		fmt.Println("  make lint     # Run linter")
	}
}
