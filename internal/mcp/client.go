package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/launch01/mission-control/internal/logging"
)

// Client represents an MCP client
type Client struct {
	baseURL    string
	httpClient *http.Client
	authMode   string
	token      string
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// NewClient creates a new MCP client
func NewClient(baseURL, authMode string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authMode: authMode,
	}
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// Call makes a JSON-RPC call to the MCP server
func (c *Client) Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  method,
		Params:  params,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	logging.Debug("MCP Request: %s %s", method, string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication based on mode
	if c.token != "" {
		if c.authMode == "header" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		} else if c.authMode == "context" {
			// For context mode, we'd add token to params
			// This is server-dependent
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	logging.Debug("MCP Response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP server returned status %d: %s", resp.StatusCode, string(body))
	}

	var response JSONRPCResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", response.Error.Code, response.Error.Message)
	}

	return response.Result, nil
}

// ListTools lists all available tools
func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	result, err := c.Call(ctx, "tools/list", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	var response struct {
		Tools []Tool `json:"tools"`
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("failed to parse tools response: %w", err)
	}

	return response.Tools, nil
}

// CallTool calls a specific tool
func (c *Client) CallTool(ctx context.Context, name string, arguments map[string]interface{}) (json.RawMessage, error) {
	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}

	return c.Call(ctx, "tools/call", params)
}
