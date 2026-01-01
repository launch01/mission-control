#!/bin/bash

# Mission Control Setup Script
# This script helps configure your environment for mission-control

set -e

echo "======================================"
echo "Mission Control Setup"
echo "======================================"
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/mission-control" ]; then
    echo "Error: Please run this script from the mission-control directory"
    exit 1
fi

# Check if .env already exists
if [ -f ".env" ]; then
    echo "⚠️  .env file already exists!"
    read -p "Do you want to overwrite it? (y/N): " overwrite
    if [ "$overwrite" != "y" ] && [ "$overwrite" != "Y" ]; then
        echo "Keeping existing .env file. Exiting."
        exit 0
    fi
fi

# Prompt for HubSpot Client ID
echo ""
echo "Step 1: HubSpot OAuth Configuration"
echo "------------------------------------"
echo ""
echo "First, create a HubSpot OAuth app if you haven't already:"
echo "1. Go to https://developers.hubspot.com/"
echo "2. Click 'Apps' → 'Create app' → 'Public app'"
echo "3. Set redirect URL: http://127.0.0.1:8400/oauth/callback"
echo "4. Add scopes: crm.objects.contacts, companies, deals (read & write)"
echo "5. Copy your Client ID"
echo ""

read -p "Enter your HubSpot Client ID: " client_id

if [ -z "$client_id" ]; then
    echo "Error: Client ID cannot be empty"
    exit 1
fi

# Create .env file
echo ""
echo "Creating .env file..."
cat > .env << EOF
# HubSpot OAuth Configuration
HUBSPOT_CLIENT_ID=$client_id
HUBSPOT_CLIENT_SECRET=  # Not needed for PKCE
HUBSPOT_REDIRECT_URI=http://127.0.0.1:8400/oauth/callback
HUBSPOT_SCOPES=crm.objects.contacts.read crm.objects.contacts.write crm.objects.companies.read crm.objects.companies.write crm.objects.deals.read crm.objects.deals.write

# MCP Server Configuration
HUBSPOT_MCP_URL=http://127.0.0.1:3333
HUBSPOT_MCP_AUTH_MODE=header

# Debug (optional)
DEBUG=false
EOF

echo "✅ .env file created successfully!"
echo ""

# Build the binary
echo "Step 2: Building mission-control..."
echo "------------------------------------"
echo ""

if command -v go &> /dev/null; then
    echo "Building mission-control binary..."
    go build -o mission-control ./cmd/mission-control
    echo "✅ Build successful!"
else
    echo "⚠️  Go not found in PATH. Trying /opt/homebrew/bin/go..."
    if [ -x "/opt/homebrew/bin/go" ]; then
        /opt/homebrew/bin/go build -o mission-control ./cmd/mission-control
        echo "✅ Build successful!"
    else
        echo "❌ Error: Go compiler not found"
        echo "Please install Go or add it to your PATH"
        exit 1
    fi
fi

echo ""
echo "======================================"
echo "Setup Complete! 🎉"
echo "======================================"
echo ""
echo "Next steps:"
echo ""
echo "1. Load environment variables:"
echo "   source .env"
echo ""
echo "2. Make sure the HubSpot MCP server is running:"
echo "   npx @hubspot/mcp-server"
echo ""
echo "3. Authenticate with HubSpot:"
echo "   ./mission-control auth login"
echo ""
echo "4. Check authentication status:"
echo "   ./mission-control auth status"
echo ""
echo "5. Test connection:"
echo "   ./mission-control tools list"
echo ""
echo "For more information, see SETUP_GUIDE.md"
echo ""
