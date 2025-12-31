package cli

import (
	"fmt"
	"os"

	"github.com/launch01/mission-control/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	mcpURL    string
	authMode  string
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "mission-control",
	Short: "HubSpot MCP Agent CLI",
	Long:  `Mission Control is a CLI tool for interacting with HubSpot via MCP (Model Context Protocol) using OAuth 2.0 authentication.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override config with flags if provided
		if mcpURL != "" {
			cfg.MCP.URL = mcpURL
		}
		if authMode != "" {
			cfg.MCP.AuthMode = authMode
		}

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&mcpURL, "mcp-url", "", "MCP server URL (default from HUBSPOT_MCP_URL or http://127.0.0.1:3333)")
	RootCmd.PersistentFlags().StringVar(&authMode, "auth-mode", "", "Authentication mode: header or context (default: header)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
