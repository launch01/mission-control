package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// OAuth endpoints
	HubSpotAuthURL  = "https://app.hubspot.com/oauth/authorize"
	HubSpotTokenURL = "https://api.hubapi.com/oauth/v1/token"

	// Default values
	DefaultRedirectURI = "http://127.0.0.1:8400/oauth/callback"
	DefaultMCPURL      = "http://127.0.0.1:3333"
	DefaultCallbackPort = "8400"
)

// Config holds application configuration
type Config struct {
	HubSpot HubSpotConfig
	MCP     MCPConfig
}

// HubSpotConfig holds HubSpot OAuth configuration
type HubSpotConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       string
}

// MCPConfig holds MCP server configuration
type MCPConfig struct {
	URL      string
	AuthMode string // "header" or "context"
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	viper.SetEnvPrefix("HUBSPOT")
	viper.AutomaticEnv()

	cfg := &Config{
		HubSpot: HubSpotConfig{
			ClientID:     getEnvOrDefault("HUBSPOT_CLIENT_ID", ""),
			ClientSecret: getEnvOrDefault("HUBSPOT_CLIENT_SECRET", ""),
			RedirectURI:  getEnvOrDefault("HUBSPOT_REDIRECT_URI", DefaultRedirectURI),
			Scopes:       getEnvOrDefault("HUBSPOT_SCOPES", "crm.objects.contacts.read crm.objects.contacts.write crm.objects.companies.read crm.objects.companies.write crm.objects.deals.read crm.objects.deals.write"),
		},
		MCP: MCPConfig{
			URL:      getEnvOrDefault("HUBSPOT_MCP_URL", DefaultMCPURL),
			AuthMode: getEnvOrDefault("HUBSPOT_MCP_AUTH_MODE", "header"),
		},
	}

	if cfg.HubSpot.ClientID == "" {
		return nil, fmt.Errorf("HUBSPOT_CLIENT_ID is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
