package cli

import (
	"context"
	"fmt"

	"github.com/launch01/mission-control/internal/agent"
	"github.com/launch01/mission-control/internal/logging"
	"github.com/launch01/mission-control/internal/oauth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with HubSpot OAuth",
	RunE: func(cmd *cobra.Command, args []string) error {
		flow, err := oauth.NewAuthFlow(cfg)
		if err != nil {
			return fmt.Errorf("failed to create auth flow: %w", err)
		}

		logging.Info("Starting OAuth login flow...")
		ctx := context.Background()
		if err := flow.Login(ctx); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		logging.Info("Successfully authenticated!")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		ag, err := agent.NewAgent(cfg)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		status, err := ag.GetAuthStatus()
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		if !status.Authenticated {
			fmt.Println("Status: Not authenticated")
			fmt.Println("Run 'mission-control auth login' to authenticate")
			return nil
		}

		fmt.Println("Status: Authenticated")
		fmt.Printf("Message: %s\n", status.Message)
		fmt.Printf("Token expires at: %s\n", status.ExpiresAt.Format("2006-01-02 15:04:05"))
		if status.IsExpired {
			fmt.Println("Note: Token is expired and will be refreshed on next use")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
}
