# Mission Control Setup Guide

## Step 1: Create HubSpot OAuth App

### 1.1 Navigate to HubSpot Developer Portal

Go to: https://developers.hubspot.com/

### 1.2 Create a New App

1. Click on **"Apps"** in the top navigation
2. Click **"Create app"** button
3. Choose **"Public app"** (this is important - NOT a Private App)

### 1.3 Configure Your App

**Basic Information:**
- **App Name**: `Mission Control CLI` (or any name you prefer)
- **Description**: `CLI tool for HubSpot CRM operations via MCP`

**Auth Tab - OAuth Settings:**

1. **Redirect URLs** - Add this URL:
   ```
   http://127.0.0.1:8400/oauth/callback
   ```

2. **Scopes** - Select the following scopes:
   - ✅ `crm.objects.contacts.read`
   - ✅ `crm.objects.contacts.write`
   - ✅ `crm.objects.companies.read`
   - ✅ `crm.objects.companies.write`
   - ✅ `crm.objects.deals.read`
   - ✅ `crm.objects.deals.write`

3. Click **"Create app"** or **"Save"**

### 1.4 Get Your Client ID

After creating the app:
1. Go to the **"Auth"** tab
2. Find **"Client ID"** section
3. **Copy the Client ID** - you'll need this in the next step

**Note**: You do NOT need the Client Secret for this setup, as we're using PKCE (Proof Key for Code Exchange) for security.

---

## Step 2: Configure Mission Control

### 2.1 Set Environment Variables

Open your terminal and run:

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control

# Set your HubSpot Client ID (replace with your actual Client ID)
export HUBSPOT_CLIENT_ID="your-client-id-from-step-1.4"

# Set other required variables (these are defaults, but good to set explicitly)
export HUBSPOT_REDIRECT_URI="http://127.0.0.1:8400/oauth/callback"
export HUBSPOT_SCOPES="crm.objects.contacts.read crm.objects.contacts.write crm.objects.companies.read crm.objects.companies.write crm.objects.deals.read crm.objects.deals.write"
export HUBSPOT_MCP_URL="http://127.0.0.1:3333"
```

### 2.2 Optional: Create .env File

For persistent configuration, create a `.env` file:

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
cat > .env << 'EOF'
# HubSpot OAuth Configuration
HUBSPOT_CLIENT_ID=your-client-id-from-step-1.4
HUBSPOT_CLIENT_SECRET=  # Not needed for PKCE
HUBSPOT_REDIRECT_URI=http://127.0.0.1:8400/oauth/callback
HUBSPOT_SCOPES=crm.objects.contacts.read crm.objects.contacts.write crm.objects.companies.read crm.objects.companies.write crm.objects.deals.read crm.objects.deals.write

# MCP Server Configuration
HUBSPOT_MCP_URL=http://127.0.0.1:3333
HUBSPOT_MCP_AUTH_MODE=header

# Debug (optional)
DEBUG=false
EOF
```

Then load it:
```bash
source .env
```

### 2.3 Build Mission Control

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
make build
```

Or build manually:
```bash
go build -o mission-control ./cmd/mission-control
```

---

## What's Next?

After completing steps 1 and 2, you'll need to:

1. **Set up the HubSpot MCP Server** (requires Node.js)
2. **Run authentication**: `./mission-control auth login`
3. **Test the connection**: `./mission-control auth status`

---

## Troubleshooting

### Error: "HUBSPOT_CLIENT_ID is required"
- Make sure you've exported the environment variable
- Verify with: `echo $HUBSPOT_CLIENT_ID`

### Error: "Invalid redirect URI"
- Check that the redirect URI in your HubSpot app matches exactly: `http://127.0.0.1:8400/oauth/callback`
- No trailing slash, correct port (8400)

### Can't find mission-control command
- Make sure you're in the correct directory: `cd /Users/ericolden/workspace_launch/claude/mission-control`
- Run with `./mission-control` (don't forget the `./`)
- Or add to PATH: `export PATH=$PATH:$(pwd)`

---

## Quick Reference Commands

```bash
# Check if build was successful
./mission-control --version

# View help
./mission-control --help

# Authentication commands
./mission-control auth login
./mission-control auth status

# List MCP tools
./mission-control tools list

# HubSpot operations
./mission-control hubspot contacts search --email test@example.com
./mission-control hubspot deals create --name "New Deal" --amount 50000
```
