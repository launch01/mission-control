package ga

import (
	"encoding/json"
	"fmt"

	"github.com/launch01/mission-control/ga-agent/mcp"
)

// Agent represents a Google Analytics MCP agent
type Agent struct {
	client *mcp.Client
}

// Tool represents an MCP tool
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewAgent creates a new GA agent
func NewAgent(accessToken string) (*Agent, error) {
	// Create MCP client that connects to GA MCP server
	client, err := mcp.NewClient(
		"npx",
		[]string{"-y", "@google-analytics/mcp-server"},
		map[string]string{
			"GA_ACCESS_TOKEN": accessToken,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}

	// Initialize the connection
	if err := client.Initialize(map[string]interface{}{"name": "ga-agent", "version": "1.0.0"}); err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	return &Agent{
		client: client,
	}, nil
}

// Close closes the agent connection
func (a *Agent) Close() error {
	return a.client.Close()
}

// ListAvailableTools lists all available MCP tools
func (a *Agent) ListAvailableTools() ([]Tool, error) {
	resp, err := a.client.ListTools()
	if err != nil {
		return nil, err
	}

	var result struct {
		Tools []Tool `json:"tools"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse tools: %w", err)
	}

	return result.Tools, nil
}

