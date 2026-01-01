package main

import (
"bufio"
"encoding/json"
"fmt"
"io"
"os"

"github.com/launch01/hubspot-mcp-agent/hubspot"
)

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
JSONRPC string          `json:"jsonrpc"`
ID      interface{}     `json:"id"`
Method  string          `json:"method"`
Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents an outgoing MCP response
type MCPResponse struct {
JSONRPC string      `json:"jsonrpc"`
ID      interface{} `json:"id"`
Result  interface{} `json:"result,omitempty"`
Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
Code    int         `json:"code"`
Message string      `json:"message"`
Data    interface{} `json:"data,omitempty"`
}

// ServerInfo represents server information
type ServerInfo struct {
Name    string `json:"name"`
Version string `json:"version"`
}

// Tool represents an MCP tool definition
type Tool struct {
Name        string                 `json:"name"`
Description string                 `json:"description"`
InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPServer wraps the HubSpot agent as an MCP server
type MCPServer struct {
agent  *hubspot.Agent
stdin  *bufio.Scanner
stdout io.Writer
}

// NewMCPServer creates a new MCP server
func NewMCPServer(accessToken string) (*MCPServer, error) {
agent, err := hubspot.NewAgent(accessToken)
if err != nil {
return nil, fmt.Errorf("failed to create HubSpot agent: %w", err)
}

return &MCPServer{
agent:  agent,
stdin:  bufio.NewScanner(os.Stdin),
stdout: os.Stdout,
}, nil
}

// sendResponse sends a JSON-RPC response
func (s *MCPServer) sendResponse(id interface{}, result interface{}, err *MCPError) {
response := MCPResponse{
JSONRPC: "2.0",
ID:      id,
Result:  result,
Error:   err,
}

jsonData, _ := json.Marshal(response)
fmt.Fprintf(s.stdout, "%s\n", jsonData)
}

// handleInitialize handles the initialize request
func (s *MCPServer) handleInitialize(id interface{}, params json.RawMessage) {
result := map[string]interface{}{
"protocolVersion": "2024-11-05",
"capabilities": map[string]interface{}{
"tools": map[string]interface{}{},
},
"serverInfo": ServerInfo{
Name:    "mission-control-hubspot",
Version: "1.0.0",
},
}
s.sendResponse(id, result, nil)
}

// handleToolsList handles the tools/list request
func (s *MCPServer) handleToolsList(id interface{}) {
tools := []Tool{
{
Name:        "hubspot_search_contacts",
Description: "Search for contacts in HubSpot CRM. Supports filtering by email and other properties.",
InputSchema: map[string]interface{}{
"type": "object",
"properties": map[string]interface{}{
"query": map[string]interface{}{
"type":        "string",
"description": "Search query (e.g., email address or name)",
},
"limit": map[string]interface{}{
"type":        "number",
"description": "Maximum number of results to return",
"default":     10,
},
},
},
},
{
Name:        "hubspot_create_contact",
Description: "Create a new contact in HubSpot CRM",
InputSchema: map[string]interface{}{
"type": "object",
"properties": map[string]interface{}{
"email": map[string]interface{}{
"type":        "string",
"description": "Contact email address",
},
"firstname": map[string]interface{}{
"type":        "string",
"description": "First name",
},
"lastname": map[string]interface{}{
"type":        "string",
"description": "Last name",
},
"phone": map[string]interface{}{
"type":        "string",
"description": "Phone number",
},
"company": map[string]interface{}{
"type":        "string",
"description": "Company name",
},
},
"required": []string{"email"},
},
},
{
Name:        "hubspot_search_companies",
Description: "Search for companies in HubSpot CRM",
InputSchema: map[string]interface{}{
"type": "object",
"properties": map[string]interface{}{
"query": map[string]interface{}{
"type":        "string",
"description": "Search query (e.g., company name or domain)",
},
"limit": map[string]interface{}{
"type":        "number",
"description": "Maximum number of results to return",
"default":     10,
},
},
},
},
{
Name:        "hubspot_search_deals",
Description: "Search for deals in HubSpot CRM",
InputSchema: map[string]interface{}{
"type": "object",
"properties": map[string]interface{}{
"query": map[string]interface{}{
"type":        "string",
"description": "Search query (e.g., deal name)",
},
"limit": map[string]interface{}{
"type":        "number",
"description": "Maximum number of results to return",
"default":     10,
},
},
},
},
{
Name:        "hubspot_get_user_details",
Description: "Get details about the authenticated HubSpot user and account",
InputSchema: map[string]interface{}{
"type": "object",
},
},
}

result := map[string]interface{}{
"tools": tools,
}
s.sendResponse(id, result, nil)
}

// handleToolsCall handles the tools/call request
func (s *MCPServer) handleToolsCall(id interface{}, params json.RawMessage) {
var request struct {
Name      string                 `json:"name"`
Arguments map[string]interface{} `json:"arguments"`
}

if err := json.Unmarshal(params, &request); err != nil {
s.sendResponse(id, nil, &MCPError{
Code:    -32602,
Message: "Invalid params",
})
return
}

var result interface{}
var err error

switch request.Name {
case "hubspot_search_contacts":
query := ""
limit := 10
if q, ok := request.Arguments["query"].(string); ok {
query = q
}
if l, ok := request.Arguments["limit"].(float64); ok {
limit = int(l)
}
result, err = s.agent.SearchContacts(query, limit)

case "hubspot_create_contact":
result, err = s.agent.CreateContact(request.Arguments)

case "hubspot_search_companies":
query := ""
limit := 10
if q, ok := request.Arguments["query"].(string); ok {
query = q
}
if l, ok := request.Arguments["limit"].(float64); ok {
limit = int(l)
}
result, err = s.agent.SearchCompanies(query, limit)

case "hubspot_search_deals":
query := ""
limit := 10
if q, ok := request.Arguments["query"].(string); ok {
query = q
}
if l, ok := request.Arguments["limit"].(float64); ok {
limit = int(l)
}
result, err = s.agent.SearchDeals(query, limit)

case "hubspot_get_user_details":
result, err = s.agent.GetUserDetails()

default:
s.sendResponse(id, nil, &MCPError{
Code:    -32601,
Message: fmt.Sprintf("Unknown tool: %s", request.Name),
})
return
}

if err != nil {
s.sendResponse(id, nil, &MCPError{
Code:    -32603,
Message: err.Error(),
})
return
}

// Parse the JSON result from HubSpot
var parsedResult interface{}
if resultBytes, ok := result.(json.RawMessage); ok {
json.Unmarshal(resultBytes, &parsedResult)
}

s.sendResponse(id, map[string]interface{}{
"content": []map[string]interface{}{
{
"type": "text",
"text": string(result.(json.RawMessage)),
},
},
}, nil)
}

// Start starts the MCP server
func (s *MCPServer) Start() error {
fmt.Fprintln(os.Stderr, "Mission Control HubSpot MCP Server started")

for s.stdin.Scan() {
line := s.stdin.Bytes()

var request MCPRequest
if err := json.Unmarshal(line, &request); err != nil {
fmt.Fprintf(os.Stderr, "Error parsing request: %v\n", err)
continue
}

switch request.Method {
case "initialize":
s.handleInitialize(request.ID, request.Params)
case "tools/list":
s.handleToolsList(request.ID)
case "tools/call":
s.handleToolsCall(request.ID, request.Params)
default:
s.sendResponse(request.ID, nil, &MCPError{
Code:    -32601,
Message: fmt.Sprintf("Method not found: %s", request.Method),
})
}
}

return s.stdin.Err()
}

func main() {
accessToken := os.Getenv("HUBSPOT_ACCESS_TOKEN")
if accessToken == "" {
fmt.Fprintln(os.Stderr, "Error: HUBSPOT_ACCESS_TOKEN environment variable required")
os.Exit(1)
}

server, err := NewMCPServer(accessToken)
if err != nil {
fmt.Fprintf(os.Stderr, "Error creating MCP server: %v\n", err)
os.Exit(1)
}

if err := server.Start(); err != nil {
fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n", err)
os.Exit(1)
}
}
