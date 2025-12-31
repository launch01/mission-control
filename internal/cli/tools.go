package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/launch01/mission-control/internal/agent"
	"github.com/spf13/cobra"
)

var (
	toolName  string
	toolInput string
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "MCP tools commands",
}

var listToolsCmd = &cobra.Command{
	Use:   "list",
	Short: "List available MCP tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		ag, err := agent.NewAgent(cfg)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		ctx := context.Background()
		tools, err := ag.ListTools(ctx)
		if err != nil {
			return fmt.Errorf("failed to list tools: %w", err)
		}

		fmt.Printf("Available tools (%d):\n\n", len(tools))
		for _, tool := range tools {
			fmt.Printf("Name: %s\n", tool.Name)
			fmt.Printf("Description: %s\n", tool.Description)
			if tool.InputSchema != nil {
				schema, _ := json.MarshalIndent(tool.InputSchema, "  ", "  ")
				fmt.Printf("Input Schema:\n  %s\n", string(schema))
			}
			fmt.Println()
		}

		return nil
	},
}

var callToolCmd = &cobra.Command{
	Use:   "call",
	Short: "Call an MCP tool",
	Example: `  mission-control tools call --name search_contacts --input '{"query": "example.com", "limit": 10}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if toolName == "" {
			return fmt.Errorf("--name is required")
		}
		if toolInput == "" {
			return fmt.Errorf("--input is required")
		}

		var inputArgs map[string]interface{}
		if err := json.Unmarshal([]byte(toolInput), &inputArgs); err != nil {
			return fmt.Errorf("invalid JSON input: %w", err)
		}

		ag, err := agent.NewAgent(cfg)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		ctx := context.Background()
		result, err := ag.CallTool(ctx, toolName, inputArgs)
		if err != nil {
			return fmt.Errorf("tool call failed: %w", err)
		}

		// Pretty print result
		var formatted interface{}
		json.Unmarshal(result, &formatted)
		output, _ := json.MarshalIndent(formatted, "", "  ")
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(toolsCmd)
	toolsCmd.AddCommand(listToolsCmd)
	toolsCmd.AddCommand(callToolCmd)

	callToolCmd.Flags().StringVarP(&toolName, "name", "n", "", "Tool name (required)")
	callToolCmd.Flags().StringVarP(&toolInput, "input", "i", "", "Tool input as JSON (required)")
}
