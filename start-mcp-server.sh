#!/bin/bash

# Start HubSpot MCP Server
# This script starts the MCP server with your HubSpot token

set -e

cd /Users/ericolden/workspace_launch/claude/mission-control

# Load environment variables
if [ ! -f ".env" ]; then
    echo "❌ Error: .env file not found"
    echo "Please create a .env file with your HUBSPOT_ACCESS_TOKEN"
    exit 1
fi

source .env

# Check if token is set
if [ -z "$HUBSPOT_ACCESS_TOKEN" ]; then
    echo "❌ Error: HUBSPOT_ACCESS_TOKEN not set in .env"
    exit 1
fi

echo "======================================"
echo "Starting HubSpot MCP Server"
echo "======================================"
echo ""
echo "Server will run on: http://127.0.0.1:3333"
echo "Press Ctrl+C to stop"
echo ""
echo "Keep this terminal window open!"
echo "======================================"
echo ""

# Start the MCP server with proper PATH
export PATH="/opt/homebrew/bin:$PATH"
npx @hubspot/mcp-server
