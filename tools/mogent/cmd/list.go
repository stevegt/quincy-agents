package cmd

import (
	"fmt"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/scope"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available modules",
	Long:  "Show all available modules and which are active for current scopes",
	RunE: func(cmd *cobra.Command, args []string) error {
		activeOnly, _ := cmd.Flags().GetBool("active")

		config, err := toml.Parse("AGENTS.toml")
		if err != nil {
			return fmt.Errorf("failed to parse AGENTS.toml: %w", err)
		}

		resolver := scope.NewResolver(config.Activate.Scopes)

		for catName, cat := range config.Category {
			for _, mod := range cat.Modules {
				active := resolver.Matches(mod.Tags)
				if activeOnly && !active {
					continue
				}

				status := "inactive"
				if active {
					status = "active"
				}

				fmt.Printf("[%s] %s/%s (%s)\n", status, catName, mod.Name, mod.Source)
			}
		}

		return nil
	},
}

func init() {
	listCmd.Flags().Bool("active", false, "Show only active modules")
}
