package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available modules",
	Long:  "Show all available modules and which are active for current scopes",
	RunE: func(cmd *cobra.Command, args []string) error {
		activeOnly, _ := cmd.Flags().GetBool("active")
		if activeOnly {
			fmt.Println("Active modules:")
		} else {
			fmt.Println("All modules:")
		}
		return nil
	},
}

func init() {
	listCmd.Flags().Bool("active", false, "Show only active modules")
}
