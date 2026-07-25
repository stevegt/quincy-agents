package cmd

import (
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "mogent",
	Short:   "Modular agent prompt manager",
	Long:    "Assemble AGENTS.md files from modular, scoped templates",
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(tagsCmd)
}
