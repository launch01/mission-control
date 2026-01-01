# HubSpot MCP Agent (Go)

A Go agent that communicates with the HubSpot MCP (Model Context Protocol) server to interact with HubSpot CRM data.

## Features

- ✅ MCP client library for Go
- ✅ Connect to HubSpot MCP server
- ✅ List available tools
- ✅ Search and manage contacts
- ✅ Search and manage companies
- ✅ Search and manage deals
- ✅ Create and update CRM objects

## Prerequisites

- Go 1.21 or later
- Node.js and npm (for HubSpot MCP server)
- HubSpot MCP server installed: `npm install -g @hubspot/mcp-server`
- HubSpot access token (from a Private App)

## Installation

1. Clone or download this repository

2. Install dependencies:
```bash
cd hubspot-mcp-agent
go mod download
```

## Getting Your HubSpot Access Token

1. Log into your HubSpot account
2. Go to **Settings** → **Integrations** → **Private Apps**
3. Click **Create a private app**
4. Give it a name (e.g., "MCP Agent")
5. Select the scopes you need:
   - `crm.objects.contacts.read` and `crm.objects.contacts.write`
   - `crm.objects.companies.read` and `crm.objects.companies.write`
   - `crm.objects.deals.read` and `crm.objects.deals.write`
6. Click **Create app**
7. Copy the access token

## Usage

### Set your HubSpot access token:

```bash
export HUBSPOT_ACCESS_TOKEN="your-access-token-here"
```

### Run the agent:

```bash
go run main.go
```

### Build and run:

```bash
go build -o hubspot-agent
./hubspot-agent
```

## Example Code

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/launch01/hubspot-mcp-agent/hubspot"
)

func main() {
	// Create agent
	agent, err := hubspot.NewAgent(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()

	// List available tools
	tools, _ := agent.ListAvailableTools()
	for _, tool := range tools {
		fmt.Printf("%s: %s\n", tool.Name, tool.Description)
	}

	// Search contacts
	contacts, _ := agent.SearchContacts("@example.com", 10)
	fmt.Println(string(contacts))

	// Create a contact
	newContact, _ := agent.CreateContact(map[string]interface{}{
		"email":     "jane@example.com",
		"firstname": "Jane",
		"lastname":  "Smith",
	})
	fmt.Println(string(newContact))

	// Create a company
	newCompany, _ := agent.CreateCompany(map[string]interface{}{
		"name":   "Acme Corp",
		"domain": "acme.com",
	})
	fmt.Println(string(newCompany))

	// Create a deal
	newDeal, _ := agent.CreateDeal(map[string]interface{}{
		"dealname": "Q1 Partnership",
		"amount":   "100000",
	})
	fmt.Println(string(newDeal))
}
```

## Project Structure

```
hubspot-mcp-agent/
├── go.mod                  # Go module definition
├── main.go                 # Example usage
├── README.md              # This file
├── mcp/
│   └── client.go          # Generic MCP client library
└── hubspot/
    └── agent.go           # HubSpot-specific agent
```

## Available Methods

### Agent Methods

- `NewAgent(accessToken string)` - Create a new HubSpot agent
- `ListAvailableTools()` - Get all available MCP tools
- `SearchContacts(query, limit)` - Search for contacts
- `GetContact(contactID)` - Get a specific contact
- `CreateContact(properties)` - Create a new contact
- `UpdateContact(contactID, properties)` - Update a contact
- `SearchCompanies(query, limit)` - Search for companies
- `CreateCompany(properties)` - Create a new company
- `SearchDeals(query, limit)` - Search for deals
- `CreateDeal(properties)` - Create a new deal
- `Close()` - Close the agent connection

## How It Works

1. The Go agent spawns the HubSpot MCP server as a subprocess using `npx @hubspot/mcp-server`
2. Communication happens over stdio using JSON-RPC 2.0 protocol
3. The MCP client library handles:
   - Process management
   - JSON-RPC request/response handling
   - Bidirectional communication over stdin/stdout
4. The HubSpot agent provides high-level methods for common CRM operations

## Error Handling

All methods return errors that should be checked:

```go
contacts, err := agent.SearchContacts("", 10)
if err != nil {
	log.Printf("Error searching contacts: %v", err)
	return
}
```

## License

MIT

## Contributing

Contributions welcome! Please submit issues and pull requests.
