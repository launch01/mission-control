# Mission Control - HubSpot MCP Agent

A secure CLI tool for interacting with HubSpot CRM via the Model Context Protocol (MCP) using OAuth 2.0 with PKCE authentication.

## Features

- ✅ **OAuth 2.0 with PKCE** - Secure authentication without client secrets
- ✅ **Token Management** - Automatic refresh with OS keychain storage (macOS/Linux) or encrypted file fallback
- ✅ **MCP Integration** - Call any HubSpot MCP tool via JSON-RPC 2.0
- ✅ **CLI Commands** - Intuitive commands for auth, tools, and HubSpot operations
- ✅ **Secure by Default** - Never prints tokens, restricted file permissions, automatic redaction
- ✅ **Retry & Backoff** - Built-in timeout and error handling

## Prerequisites

- Go 1.21 or later
- HubSpot account with OAuth app configured
- Local MCP server running (default: http://127.0.0.1:3333)

## Installation

### From Source

```bash
git clone https://github.com/launch01/mission-control.git
cd mission-control
make build
make install  # Optional: install to $GOPATH/bin
```

### Quick Build

```bash
go build -o mission-control ./cmd/mission-control
```

## HubSpot OAuth App Setup

### 1. Create a Public App in HubSpot

1. Log into your HubSpot account
2. Navigate to **Settings** → **Integrations** → **Private Apps** or use the [HubSpot Developer Portal](https://developers.hubspot.com/)
3. Create a new **Public App** (not Private App)
4. Configure the app:
   - **App Name**: Mission Control CLI
   - **Redirect URL**: `http://127.0.0.1:8400/oauth/callback`
   - **Scopes** (minimum required):
     - `crm.objects.contacts.read`
     - `crm.objects.contacts.write`
     - `crm.objects.companies.read`
     - `crm.objects.companies.write`
     - `crm.objects.deals.read`
     - `crm.objects.deals.write`
5. Save and copy your **Client ID**
6. (Optional) Copy **Client Secret** if using confidential client mode

### 2. Configure Environment Variables

```bash
cp .env.example .env
# Edit .env with your values
```

Required variables:
```bash
HUBSPOT_CLIENT_ID=your-client-id-here
HUBSPOT_CLIENT_SECRET=  # Optional for PKCE
HUBSPOT_REDIRECT_URI=http://127.0.0.1:8400/oauth/callback
HUBSPOT_SCOPES=crm.objects.contacts.read crm.objects.contacts.write crm.objects.companies.read crm.objects.companies.write crm.objects.deals.read crm.objects.deals.write
HUBSPOT_MCP_URL=http://127.0.0.1:3333
```

Load environment:
```bash
source .env
# Or use export for each variable
export HUBSPOT_CLIENT_ID=your-client-id
```

## Running the Local MCP Server

The MCP server must be running before using mission-control.

### Option 1: HubSpot MCP Server (if available)

```bash
npm install -g @hubspot/mcp-server
HUBSPOT_ACCESS_TOKEN=temp npx @hubspot/mcp-server
```

Note: The mission-control CLI will provide the actual access token via OAuth.

### Option 2: Custom MCP Server

Implement a JSON-RPC 2.0 server that exposes HubSpot tools at `http://127.0.0.1:3333`. See [MCP Specification](https://github.com/anthropics/mcp) for details.

## Usage

### Authentication

#### Login

```bash
mission-control auth login
```

This will:
1. Generate PKCE code verifier and challenge
2. Open your browser to HubSpot authorization page
3. Start a local callback server on port 8400
4. Exchange authorization code for tokens
5. Store tokens securely in OS keychain (or encrypted file)

#### Check Status

```bash
mission-control auth status
```

Output:
```
Status: Authenticated
Message: Authenticated
Token expires at: 2025-12-31 15:30:00
```

### MCP Tools

#### List Available Tools

```bash
mission-control tools list
```

#### Call a Tool

```bash
mission-control tools call --name search_contacts --input '{"query": "example.com", "limit": 10}'
```

Example output:
```json
{
  "results": [
    {
      "id": "12345",
      "properties": {
        "email": "john@example.com",
        "firstname": "John",
        "lastname": "Doe"
      }
    }
  ]
}
```

### HubSpot Convenience Commands

#### Search Contacts

```bash
mission-control hubspot contacts search --email john@example.com
```

#### Create Deal

```bash
mission-control hubspot deals create --name "New Partnership" --amount 50000
```

### Configuration Flags

Override environment variables with flags:

```bash
mission-control --mcp-url http://localhost:3333 --auth-mode header tools list
```

Available flags:
- `--mcp-url`: MCP server URL (default: http://127.0.0.1:3333)
- `--auth-mode`: Authentication mode - `header` (default) or `context`

## Architecture

```
┌─────────────────────────────────────┐
│      mission-control CLI            │
├─────────────────────────────────────┤
│  Commands (auth, tools, hubspot)    │
├─────────────────────────────────────┤
│  Agent (orchestrates OAuth + MCP)   │
├──────────────┬──────────────────────┤
│ OAuth Flow   │   MCP Client         │
│ (PKCE)       │   (JSON-RPC 2.0)     │
├──────────────┼──────────────────────┤
│ Token Storage (Keychain/File)       │
└─────────────────────────────────────┘
         │                    │
         ▼                    ▼
  HubSpot OAuth        Local MCP Server
  Authorization              │
                             ▼
                      HubSpot REST API
```

## Token Storage

Tokens are stored securely using:

1. **OS Keychain** (preferred):
   - macOS: Keychain Access
   - Linux: Secret Service (gnome-keyring, kwallet)

2. **File Fallback** (if keychain unavailable):
   - Location: `~/.config/mission-control/token.json`
   - Permissions: `0600` (owner read/write only)
   - ⚠️ Warning displayed on first use

## Troubleshooting

### Port 8400 Already in Use

Change the redirect URI:
```bash
export HUBSPOT_REDIRECT_URI=http://127.0.0.1:8401/oauth/callback
```

Update your HubSpot app's redirect URL to match.

### Invalid Redirect URI Error

Ensure the redirect URI in your `.env` matches exactly what's configured in your HubSpot app (including port and path).

### Token Refresh Failed

If refresh fails:
```bash
mission-control auth login  # Re-authenticate
```

### MCP Server Connection Failed

1. Verify MCP server is running:
   ```bash
   curl http://127.0.0.1:3333
   ```

2. Check MCP URL:
   ```bash
   mission-control --mcp-url http://localhost:3333 tools list
   ```

### Debug Mode

Enable debug logging:
```bash
export DEBUG=true
mission-control tools list
```

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Test Coverage

```bash
make test-coverage
```

### Lint

```bash
make lint
```

Requires `golangci-lint`:
```bash
brew install golangci-lint  # macOS
```

### Clean

```bash
make clean
```

## Testing

### Unit Tests

```bash
go test ./...
```

### Integration Testing

Set up test environment:
```bash
export HUBSPOT_CLIENT_ID=test-client-id
export HUBSPOT_MCP_URL=http://localhost:3333
```

Run integration tests (requires running MCP server):
```bash
make test
```

## Security

### Best Practices

- ✅ Tokens never printed to stdout/logs
- ✅ Sensitive data redacted in debug logs
- ✅ File permissions restricted to 0600
- ✅ OAuth state parameter prevents CSRF
- ✅ PKCE prevents authorization code interception
- ✅ Automatic token refresh before expiry
- ✅ Secure token storage (keychain preferred)

### Sensitive Data Handling

All access tokens are:
- Stored encrypted in OS keychain or restricted file
- Redacted in logs (shows only first 8 and last 4 characters)
- Never passed as command-line arguments
- Automatically refreshed when expired

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `make test` and `make lint` pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

- GitHub Issues: https://github.com/launch01/mission-control/issues
- HubSpot API Docs: https://developers.hubspot.com/docs/api/overview
- MCP Specification: https://github.com/anthropics/mcp

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Uses [go-keyring](https://github.com/zalando/go-keyring) for secure storage
- OAuth 2.0 implementation based on [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
