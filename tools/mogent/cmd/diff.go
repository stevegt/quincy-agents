package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [scope]",
	Short: "Compare modules across scopes",
	Long:  "Show section-aware diff between the same module in different scopes",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("Comparing local vs assembled...")
		} else {
			fmt.Printf("Comparing with scope: %s\n", args[0])
		}
		return nil
	},
}
