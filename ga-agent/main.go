package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/launch01/mission-control/ga-agent/ga"
)

func main() {
	// Get GA access token from environment
	accessToken := os.Getenv("GA_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("GA_ACCESS_TOKEN environment variable is required")
	}

	// Create GA agent
	agent, err := ga.NewAgent(accessToken)
	if err != nil {
		log.Fatalf("Failed to create GA agent: %v", err)
	}
	defer agent.Close()

	fmt.Println("GA Agent Started!")
	fmt.Println("========================================")

	// Example 1: List available tools
	fmt.Println("\n1. Listing available tools...")
	tools, err := agent.ListAvailableTools()
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
	} else {
		fmt.Printf("Found %d tools:\n", len(tools))
		for _, tool := range tools {
			fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("Agent execution completed!")
}

