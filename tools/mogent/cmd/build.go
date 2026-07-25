package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Assemble AGENTS.md from modules",
	Long:  "Parse AGENTS.toml, resolve active scopes, and assemble the final AGENTS.md",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry")
		if dryRun {
			fmt.Println("Dry run - would assemble AGENTS.md")
			return nil
		}
		fmt.Println("Building AGENTS.md...")
		return nil
	},
}

func init() {
	buildCmd.Flags().Bool("dry", false, "Preview output without writing")
}
