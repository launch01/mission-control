package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

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

// Client represents an MCP client
type Client struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	scanner   *bufio.Scanner
	responses map[string]chan *JSONRPCResponse
	mu        sync.Mutex
	closed    bool
}

// NewClient creates a new MCP client that communicates with an MCP server
func NewClient(command string, args []string, env map[string]string) (*Client, error) {
	cmd := exec.Command(command, args...)

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	client := &Client{
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		scanner:   bufio.NewScanner(stdout),
		responses: make(map[string]chan *JSONRPCResponse),
	}

	// Start reading responses
	go client.readResponses()

	// Log stderr
	go client.logStderr()

	return client, nil
}

// readResponses reads JSON-RPC responses from stdout
func (c *Client) readResponses() {
	for c.scanner.Scan() {
		line := c.scanner.Bytes()

		var response JSONRPCResponse
		if err := json.Unmarshal(line, &response); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse response: %v\n", err)
			continue
		}

		c.mu.Lock()
		if ch, ok := c.responses[response.ID]; ok {
			ch <- &response
			delete(c.responses, response.ID)
		}
		c.mu.Unlock()
	}
}

// logStderr logs stderr output
func (c *Client) logStderr() {
	scanner := bufio.NewScanner(c.stderr)
	for scanner.Scan() {
		fmt.Fprintf(os.Stderr, "[MCP Server] %s\n", scanner.Text())
	}
}

// Call sends a JSON-RPC request and waits for the response
func (c *Client) Call(method string, params interface{}) (json.RawMessage, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	id := uuid.New().String()
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	// Register response channel
	responseChan := make(chan *JSONRPCResponse, 1)
	c.mu.Lock()
	c.responses[id] = responseChan
	c.mu.Unlock()

	// Send request
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	requestBytes = append(requestBytes, '\n')
	if _, err := c.stdin.Write(requestBytes); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	// Wait for response
	response := <-responseChan
	if response.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error %d: %s", response.Error.Code, response.Error.Message)
	}

	return response.Result, nil
}

// Initialize sends the initialize request to the MCP server
func (c *Client) Initialize(clientInfo map[string]interface{}) error {
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      clientInfo,
	}

	result, err := c.Call("initialize", params)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	fmt.Printf("Initialized: %s\n", string(result))
	return nil
}

// ListTools returns the list of available tools from the MCP server
func (c *Client) ListTools() (json.RawMessage, error) {
	return c.Call("tools/list", map[string]interface{}{})
}

// CallTool calls a specific tool on the MCP server
func (c *Client) CallTool(name string, arguments map[string]interface{}) (json.RawMessage, error) {
	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}
	return c.Call("tools/call", params)
}

// Close closes the MCP client
func (c *Client) Close() error {
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()

	c.stdin.Close()
	c.stdout.Close()
	c.stderr.Close()

	return c.cmd.Wait()
}
