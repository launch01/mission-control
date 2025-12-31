package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMCPClientCall(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Return mock response
		response := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      "test-id",
			Result:  json.RawMessage(`{"status": "success"}`),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "header")
	client.SetToken("test-token")

	result, err := client.Call(context.Background(), "test/method", map[string]interface{}{"test": "value"})
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if data["status"] != "success" {
		t.Errorf("Expected status=success, got %v", data["status"])
	}
}

func TestMCPClientAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token-123" {
			t.Errorf("Expected Authorization header 'Bearer test-token-123', got %s", auth)
		}

		response := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      "test-id",
			Result:  json.RawMessage(`{}`),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "header")
	client.SetToken("test-token-123")

	_, err := client.Call(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}
}

func TestMCPClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      "test-id",
			Error: &JSONRPCError{
				Code:    -32600,
				Message: "Invalid request",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "header")

	_, err := client.Call(context.Background(), "test", nil)
	if err == nil {
		t.Error("Expected error from MCP call")
	}
}
