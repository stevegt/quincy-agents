package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mogent",
	Short: "Modular agent prompt manager",
	Long:  "Assemble AGENTS.md files from modular, scoped templates",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(listCmd)
}
