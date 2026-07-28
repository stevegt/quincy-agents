package cmd

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/diff"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <module>",
	Short: "Compare modules across scopes",
	Long:  "Show which sections of a module are included/excluded by active scopes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := toml.Parse("AGENTS.toml")
		if err != nil {
			return fmt.Errorf("failed to parse AGENTS.toml: %w", err)
		}

		moduleName := args[0]
		moduleSource := moduleName
		if !strings.HasSuffix(moduleSource, ".md") {
			moduleSource += ".md"
		}
		path := config.ResolveModulePath(moduleSource)

		mod, err := module.Parse(path)
		if err != nil {
			return fmt.Errorf("failed to parse module %s: %w", moduleName, err)
		}

		result := diff.AnalyzeModule(mod, config.Activate.Scopes)
		fmt.Print(diff.FormatAnalysis(result, moduleName))

		return nil
	},
}
