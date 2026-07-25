package cmd

import (
	"fmt"
	"strings"

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
		tagsFilter, _ := cmd.Flags().GetString("tags")
		showTags, _ := cmd.Flags().GetBool("show-tags")

		config, err := toml.Parse("AGENTS.toml")
		if err != nil {
			return fmt.Errorf("failed to parse AGENTS.toml: %w", err)
		}

		var resolver *scope.Resolver
		if tagsFilter != "" {
			scopes := strings.Split(tagsFilter, ",")
			for i := range scopes {
				scopes[i] = strings.TrimSpace(scopes[i])
			}
			resolver = scope.NewResolver(scopes)
		} else {
			resolver = scope.NewResolver(config.Activate.Scopes)
		}

		for _, catName := range config.Order.Categories {
			cat, ok := config.Category[catName]
			if !ok {
				continue
			}

			for _, mod := range cat.Modules {
				active := resolver.Matches(mod.Tags)
				if activeOnly && !active {
					continue
				}

				status := "inactive"
				if active {
					status = "active"
				}

				if showTags {
					var tags []string
					filePath := config.ResolveModulePath(mod.Source)
					fileTags := extractTagsFromFile(filePath)
					if len(fileTags) > 0 {
						tags = fileTags
					} else if len(mod.Tags) > 0 {
						tags = mod.Tags
					}

					if len(tags) > 0 {
						fmt.Printf("[%s] %s/%s (%s) {%s}\n", status, catName, mod.Name, mod.Source, strings.Join(tags, ", "))
					} else {
						fmt.Printf("[%s] %s/%s (%s)\n", status, catName, mod.Name, mod.Source)
					}
				} else {
					fmt.Printf("[%s] %s/%s (%s)\n", status, catName, mod.Name, mod.Source)
				}
			}
		}

		return nil
	},
}

func init() {
	listCmd.Flags().Bool("active", false, "Show only active modules")
	listCmd.Flags().String("tags", "", "Filter by tags (comma-separated)")
	listCmd.Flags().Bool("show-tags", false, "Show tags for each module")
}
