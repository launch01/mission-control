# Next Steps to Get Mission Control Running

## Overview
Everything is set up! You just need to:
1. Create a HubSpot OAuth app (5 minutes)
2. Run the setup script with your Client ID
3. Authenticate with HubSpot
4. Start using mission-control

---

## Step 1: Create HubSpot OAuth App

### Open HubSpot Developer Portal
Go to: **https://developers.hubspot.com/**

### Create the App
1. Click **"Apps"** in the top navigation
2. Click **"Create app"** button
3. Choose **"Public app"**

### Configure Authentication

**App Name**: `Mission Control CLI`

**Redirect URLs** (in the Auth tab):
```
http://127.0.0.1:8400/oauth/callback
```

**Scopes** (Required - select all of these):
- ✅ crm.objects.contacts.read
- ✅ crm.objects.contacts.write
- ✅ crm.objects.companies.read
- ✅ crm.objects.companies.write
- ✅ crm.objects.deals.read
- ✅ crm.objects.deals.write

### Save and Get Client ID
1. Click **"Create app"** or **"Save"**
2. Go to the **"Auth"** tab
3. Copy your **Client ID** (you'll need this in Step 2)

---

## Step 2: Run Setup Script

Open your terminal and run:

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
./quick-start.sh
```

When prompted, paste your **Client ID** from Step 1.

The script will:
- ✅ Verify Node.js is working
- ✅ Verify HubSpot MCP server is installed
- ✅ Create your `.env` file
- ✅ Build the mission-control binary

---

## Step 3: Start the HubSpot MCP Server

In a **separate terminal window**, run:

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control

# Load your environment
source .env

# Start the MCP server
PATH="/opt/homebrew/bin:$PATH" npx @hubspot/mcp-server
```

**Keep this terminal running!** The MCP server needs to stay active.

---

## Step 4: Authenticate with HubSpot

In your **original terminal** (not the one running the MCP server), run:

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control

# Load your environment
source .env

# Login to HubSpot
./mission-control auth login
```

This will:
1. Open your browser to HubSpot's authorization page
2. Ask you to approve the app
3. Save your access token securely in macOS Keychain

---

## Step 5: Test Everything Works

```bash
# Check authentication status
./mission-control auth status

# List available MCP tools
./mission-control tools list

# Try searching contacts
./mission-control hubspot contacts search --email test@example.com

# Try creating a deal
./mission-control hubspot deals create --name "Test Deal" --amount 10000
```

---

## Quick Reference

### Start MCP Server (Terminal 1)
```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
source .env
PATH="/opt/homebrew/bin:$PATH" npx @hubspot/mcp-server
```

### Use Mission Control (Terminal 2)
```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
source .env
./mission-control [command]
```

### Common Commands
```bash
# Authentication
./mission-control auth login
./mission-control auth status

# Tools
./mission-control tools list
./mission-control tools call --name search_contacts --input '{"query":"test","limit":5}'

# HubSpot shortcuts
./mission-control hubspot contacts search --email john@example.com
./mission-control hubspot deals create --name "Q1 Deal" --amount 50000
```

---

## Troubleshooting

### "HUBSPOT_CLIENT_ID is required"
Run: `source .env` to load your environment variables

### "Connection refused" or "MCP server not responding"
Make sure the MCP server is running in a separate terminal (Step 3)

### "Invalid redirect URI"
Double-check that your HubSpot app has exactly this redirect URL:
`http://127.0.0.1:8400/oauth/callback`

### Port 8400 already in use
Someone else is using that port. You can change it in your HubSpot app settings and `.env` file.

---

## What's Happening Under the Hood

```
┌─────────────────────────────────────────────────┐
│  Your Command: ./mission-control hubspot ...   │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│  mission-control (Go CLI)                       │
│  - Loads OAuth token from macOS Keychain       │
│  - Connects to local MCP server                 │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│  HubSpot MCP Server (Node.js)                   │
│  - Running on http://127.0.0.1:3333            │
│  - Receives JSON-RPC requests                   │
│  - Uses your OAuth token for auth              │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│  HubSpot REST API                               │
│  - Creates/reads/updates CRM data              │
└─────────────────────────────────────────────────┘
```

---

## You're All Set!

Once you complete these steps, you'll be able to:
- ✅ Interact with HubSpot CRM from the command line
- ✅ Use mission-control from Claude Code
- ✅ Automate HubSpot operations via scripts
- ✅ Build custom integrations

Happy coding! 🚀
