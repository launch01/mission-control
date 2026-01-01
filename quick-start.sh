#!/bin/bash

# Quick Start Script for Mission Control
# This script ensures proper PATH and runs the setup

set -e

echo "======================================"
echo "Mission Control Quick Start"
echo "======================================"
echo ""

# Ensure we're in the right directory
cd /Users/ericolden/workspace_launch/claude/mission-control

# Set up PATH to include homebrew
export PATH="/opt/homebrew/bin:$PATH"

# Verify Node.js is available
echo "Checking Node.js installation..."
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js not found even after adding /opt/homebrew/bin to PATH"
    exit 1
fi

echo "✅ Node.js version: $(node --version)"
echo "✅ npm version: $(npm --version)"
echo ""

# Check if HubSpot MCP server is installed
echo "Checking HubSpot MCP server..."
if npm list -g @hubspot/mcp-server &> /dev/null; then
    echo "✅ HubSpot MCP server is installed"
else
    echo "⚠️  HubSpot MCP server not found. Installing..."
    npm install -g @hubspot/mcp-server
    echo "✅ HubSpot MCP server installed successfully"
fi
echo ""

# Now run the main setup script
echo "Running main setup..."
echo ""
./setup.sh
