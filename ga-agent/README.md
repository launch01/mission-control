# GA Agent

A Google Analytics MCP (Model Context Protocol) agent built in Go.

## Overview

The GA Agent provides a Go-based interface to interact with Google Analytics through the MCP protocol.

## Structure

```
ga-agent/
├── main.go          # Entry point for the agent
├── ga/
│   └── agent.go     # GA agent implementation
├── mcp/
│   └── client.go    # MCP client implementation
└── README.md        # This file
```

## Usage

1. Set the `GA_ACCESS_TOKEN` environment variable:
   ```bash
   export GA_ACCESS_TOKEN=your_token_here
   ```

2. Run the agent:
   ```bash
   go run ga-agent/main.go
   ```

## Development

This agent follows the same structure as the HubSpot MCP agent in this repository, providing a consistent interface for MCP-based integrations.

