package cmd

import (
	"fmt"
	"os"
	"strings"

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
		tagsFilter, _ := cmd.Flags().GetString("tags")
		configFile, _ := cmd.Flags().GetString("config")

		if configFile == "" {
			configFile = "AGENTS.toml"
		}

		config, err := toml.Parse(configFile)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", configFile, err)
		}

		if tagsFilter != "" {
			scopes := strings.Split(tagsFilter, ",")
			for i := range scopes {
				scopes[i] = strings.TrimSpace(scopes[i])
			}
			config.Activate.Scopes = scopes
		}

		engine := assemble.NewEngine(config)
		output, err := engine.Assemble()
		if err != nil {
			return fmt.Errorf("assembly failed: %w", err)
		}
		// Intent: Refuse empty generated prompts and surface weak block choices
		// before writing AGENTS.md. Source: DI-nasot
		if err := assemble.ValidateOutput(output); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		for _, warning := range engine.Warnings() {
			fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
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
	buildCmd.Flags().String("tags", "", "Filter by tags (comma-separated, overrides AGENTS.toml)")
	buildCmd.Flags().StringP("config", "c", "AGENTS.toml", "Path to AGENTS.toml config file")
}
