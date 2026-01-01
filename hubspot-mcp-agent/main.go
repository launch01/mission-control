package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/launch01/hubspot-mcp-agent/hubspot"
)

func main() {
	// Get HubSpot access token from environment
	accessToken := os.Getenv("HUBSPOT_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("HUBSPOT_ACCESS_TOKEN environment variable is required")
	}

	// Create HubSpot agent
	agent, err := hubspot.NewAgent(accessToken)
	if err != nil {
		log.Fatalf("Failed to create HubSpot agent: %v", err)
	}
	defer agent.Close()

	fmt.Println("HubSpot MCP Agent Started!")
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

	// Example 2: Search for contacts
	fmt.Println("\n2. Searching for contacts...")
	contacts, err := agent.SearchContacts("", 5)
	if err != nil {
		log.Printf("Failed to search contacts: %v", err)
	} else {
		var prettyContacts interface{}
		json.Unmarshal(contacts, &prettyContacts)
		prettyJSON, _ := json.MarshalIndent(prettyContacts, "", "  ")
		fmt.Printf("Contacts:\n%s\n", string(prettyJSON))
	}

	// Example 3: Create a new contact
	fmt.Println("\n3. Creating a new contact...")
	newContact, err := agent.CreateContact(map[string]interface{}{
		"email":     "john.doe@example.com",
		"firstname": "John",
		"lastname":  "Doe",
		"phone":     "+1-555-0123",
		"company":   "Example Corp",
	})
	if err != nil {
		log.Printf("Failed to create contact: %v", err)
	} else {
		var prettyContact interface{}
		json.Unmarshal(newContact, &prettyContact)
		prettyJSON, _ := json.MarshalIndent(prettyContact, "", "  ")
		fmt.Printf("Created contact:\n%s\n", string(prettyJSON))
	}

	// Example 4: Search for companies
	fmt.Println("\n4. Searching for companies...")
	companies, err := agent.SearchCompanies("", 5)
	if err != nil {
		log.Printf("Failed to search companies: %v", err)
	} else {
		var prettyCompanies interface{}
		json.Unmarshal(companies, &prettyCompanies)
		prettyJSON, _ := json.MarshalIndent(prettyCompanies, "", "  ")
		fmt.Printf("Companies:\n%s\n", string(prettyJSON))
	}

	// Example 5: Create a new deal
	fmt.Println("\n5. Creating a new deal...")
	newDeal, err := agent.CreateDeal(map[string]interface{}{
		"dealname":   "New Partnership Deal",
		"amount":     "50000",
		"dealstage":  "appointmentscheduled",
		"pipeline":   "default",
		"closedate":  "2025-03-31",
	})
	if err != nil {
		log.Printf("Failed to create deal: %v", err)
	} else {
		var prettyDeal interface{}
		json.Unmarshal(newDeal, &prettyDeal)
		prettyJSON, _ := json.MarshalIndent(prettyDeal, "", "  ")
		fmt.Printf("Created deal:\n%s\n", string(prettyJSON))
	}

	fmt.Println("\n========================================")
	fmt.Println("Agent execution completed!")
}
