package cmd

import (
	"fmt"
	"os"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/assemble"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Assemble AGENTS.md from modules",
	Long:  "Parse AGENTS.toml, resolve active scopes, and assemble the final AGENTS.md",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry")

		config, err := toml.Parse("AGENTS.toml")
		if err != nil {
			return fmt.Errorf("failed to parse AGENTS.toml: %w", err)
		}

		engine := assemble.NewEngine(config)
		output, err := engine.Assemble()
		if err != nil {
			return fmt.Errorf("assembly failed: %w", err)
		}

		if dryRun {
			fmt.Println(output)
			return nil
		}

		if err := os.WriteFile(config.Output.Path, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", config.Output.Path, err)
		}

		fmt.Printf("Built %s\n", config.Output.Path)
		return nil
	},
}

func init() {
	buildCmd.Flags().Bool("dry", false, "Preview output without writing")
}
