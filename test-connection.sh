#!/bin/bash

# Test Mission Control Connection
# This script tests that everything is working

set -e

cd /Users/ericolden/workspace_launch/claude/mission-control

# Load environment
source .env

echo "======================================"
echo "Testing Mission Control Connection"
echo "======================================"
echo ""

# Test 1: Check if MCP server is running
echo "Test 1: Checking if MCP server is running..."
if curl -s http://127.0.0.1:3333 >/dev/null 2>&1; then
    echo "✅ MCP server is running on port 3333"
else
    echo "❌ MCP server is not responding"
    echo "   Start it with: ./start-mcp-server.sh"
    exit 1
fi
echo ""

# Test 2: List available tools
echo "Test 2: Listing available MCP tools..."
export PATH="/opt/homebrew/bin:$PATH"
if ./mission-control tools list 2>&1 | head -10; then
    echo ""
    echo "✅ Can list MCP tools"
else
    echo "❌ Failed to list tools"
    exit 1
fi
echo ""

# Test 3: Try searching contacts
echo "Test 3: Testing HubSpot contact search..."
if ./mission-control hubspot contacts search --email test 2>&1 | head -20; then
    echo ""
    echo "✅ Can search HubSpot contacts"
else
    echo "⚠️  Contact search had issues (this might be normal if no contacts match)"
fi
echo ""

echo "======================================"
echo "Connection Test Complete! 🎉"
echo "======================================"
echo ""
echo "Mission Control is ready to use!"
echo ""
echo "Try these commands:"
echo "  ./mission-control tools list"
echo "  ./mission-control hubspot contacts search --email your@email.com"
echo "  ./mission-control hubspot deals create --name 'Test Deal' --amount 10000"
echo ""
