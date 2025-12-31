package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/launch01/mission-control/internal/agent"
	"github.com/spf13/cobra"
)

var (
	email       string
	dealName    string
	dealAmount  string
	companyName string
)

var hubspotCmd = &cobra.Command{
	Use:   "hubspot",
	Short: "HubSpot convenience commands",
}

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Contact management commands",
}

var searchContactsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for contacts",
	Example: `  mission-control hubspot contacts search --email john@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ag, err := agent.NewAgent(cfg)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		ctx := context.Background()

		// Call MCP tool for searching contacts
		inputArgs := map[string]interface{}{
			"query": email,
			"limit": 10,
		}

		result, err := ag.CallTool(ctx, "search_contacts", inputArgs)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		var formatted interface{}
		json.Unmarshal(result, &formatted)
		output, _ := json.MarshalIndent(formatted, "", "  ")
		fmt.Println(string(output))

		return nil
	},
}

var dealsCmd = &cobra.Command{
	Use:   "deals",
	Short: "Deal management commands",
}

var createDealCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deal",
	Example: `  mission-control hubspot deals create --name "New Partnership" --amount 50000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dealName == "" {
			return fmt.Errorf("--name is required")
		}

		ag, err := agent.NewAgent(cfg)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		ctx := context.Background()

		// Call MCP tool for creating deal
		inputArgs := map[string]interface{}{
			"properties": map[string]interface{}{
				"dealname": dealName,
			},
		}

		if dealAmount != "" {
			inputArgs["properties"].(map[string]interface{})["amount"] = dealAmount
		}

		result, err := ag.CallTool(ctx, "create_deal", inputArgs)
		if err != nil {
			return fmt.Errorf("create deal failed: %w", err)
		}

		var formatted interface{}
		json.Unmarshal(result, &formatted)
		output, _ := json.MarshalIndent(formatted, "", "  ")
		fmt.Println("Deal created successfully:")
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(hubspotCmd)

	// Contacts
	hubspotCmd.AddCommand(contactsCmd)
	contactsCmd.AddCommand(searchContactsCmd)
	searchContactsCmd.Flags().StringVarP(&email, "email", "e", "", "Email to search for")

	// Deals
	hubspotCmd.AddCommand(dealsCmd)
	dealsCmd.AddCommand(createDealCmd)
	createDealCmd.Flags().StringVarP(&dealName, "name", "n", "", "Deal name (required)")
	createDealCmd.Flags().StringVarP(&dealAmount, "amount", "a", "", "Deal amount")
}
