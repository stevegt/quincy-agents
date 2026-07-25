package cmd

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags across modules",
	Long:  "Scan all modules and show which tags are available",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		filter, _ := cmd.Flags().GetString("filter")

		if configFile == "" {
			configFile = "AGENTS.toml"
		}

		config, err := toml.Parse(configFile)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", configFile, err)
		}

		tagSet := make(map[string]bool)

		for _, catName := range config.Order.Categories {
			cat, ok := config.Category[catName]
			if !ok {
				continue
			}

			for _, mod := range cat.Modules {
				filePath := config.ResolveModulePath(mod.Source)
				fileTags := extractTagsFromFile(filePath)
				for _, tag := range fileTags {
					if filter == "" || strings.HasPrefix(tag, filter) {
						tagSet[tag] = true
					}
				}

				for _, tag := range mod.Tags {
					if filter == "" || strings.HasPrefix(tag, filter) {
						tagSet[tag] = true
					}
				}
			}
		}

		if len(tagSet) == 0 {
			fmt.Println("No tags found.")
			return nil
		}

		fmt.Println("Available tags:")
		printTagTree(tagSet)

		return nil
	},
}

func init() {
	tagsCmd.Flags().String("filter", "", "Filter tags by prefix (e.g., 'team' or 'lang')")
	tagsCmd.Flags().StringP("config", "c", "AGENTS.toml", "Path to AGENTS.toml config file")
}
