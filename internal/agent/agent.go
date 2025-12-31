package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/launch01/mission-control/internal/config"
	"github.com/launch01/mission-control/internal/logging"
	"github.com/launch01/mission-control/internal/mcp"
	"github.com/launch01/mission-control/internal/oauth"
	"github.com/launch01/mission-control/internal/storage"
)

// Agent combines OAuth and MCP functionality
type Agent struct {
	cfg         *config.Config
	mcpClient   *mcp.Client
	storage     *storage.TokenStorage
	oauthFlow   *oauth.AuthFlow
}

// NewAgent creates a new agent
func NewAgent(cfg *config.Config) (*Agent, error) {
	mcpClient := mcp.NewClient(cfg.MCP.URL, cfg.MCP.AuthMode)

	store, err := storage.NewTokenStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create token storage: %w", err)
	}

	authFlow, err := oauth.NewAuthFlow(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth flow: %w", err)
	}

	return &Agent{
		cfg:       cfg,
		mcpClient: mcpClient,
		storage:   store,
		oauthFlow: authFlow,
	}, nil
}

// EnsureAuthenticated ensures we have a valid token
func (a *Agent) EnsureAuthenticated(ctx context.Context) error {
	token, err := a.storage.LoadToken()
	if err != nil {
		return fmt.Errorf("not authenticated - please run 'mission-control auth login': %w", err)
	}

	// Refresh if expired or expiring soon
	if token.IsExpired() || token.IsExpiringSoon(5*time.Minute) {
		logging.Info("Token expired or expiring soon, refreshing...")
		if err := a.oauthFlow.RefreshToken(ctx); err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		// Reload token after refresh
		token, err = a.storage.LoadToken()
		if err != nil {
			return fmt.Errorf("failed to reload token: %w", err)
		}
	}

	// Set token on MCP client
	a.mcpClient.SetToken(token.AccessToken)
	return nil
}

// ListTools lists all available MCP tools
func (a *Agent) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	if err := a.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	return a.mcpClient.ListTools(ctx)
}

// CallTool calls an MCP tool
func (a *Agent) CallTool(ctx context.Context, name string, args map[string]interface{}) (json.RawMessage, error) {
	if err := a.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	return a.mcpClient.CallTool(ctx, name, args)
}

// GetAuthStatus returns the current authentication status
func (a *Agent) GetAuthStatus() (*AuthStatus, error) {
	token, err := a.storage.LoadToken()
	if err != nil {
		return &AuthStatus{
			Authenticated: false,
			Message:       "Not authenticated",
		}, nil
	}

	status := &AuthStatus{
		Authenticated: true,
		ExpiresAt:     token.ExpiresAt,
		IsExpired:     token.IsExpired(),
	}

	if token.IsExpired() {
		status.Message = "Token expired - will be refreshed on next use"
	} else if token.IsExpiringSoon(5 * time.Minute) {
		status.Message = "Token expiring soon - will be refreshed on next use"
	} else {
		status.Message = "Authenticated"
	}

	return status, nil
}

// AuthStatus represents authentication status
type AuthStatus struct {
	Authenticated bool
	ExpiresAt     time.Time
	IsExpired     bool
	Message       string
}
