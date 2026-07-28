// Intent: Render module selections as a tree so active/inactive blocks are
// readable before the project grows a richer TUI.
// Source: DI-soviv

package cmd

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
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
		configFile, _ := cmd.Flags().GetString("config")

		if configFile == "" {
			configFile = "AGENTS.toml"
		}

		config, err := toml.Parse(configFile)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", configFile, err)
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

		selectedReferences, err := selectedReferencesByPath(config)
		if err != nil {
			return err
		}

		fmt.Println(".")
		visibleCategories := make([]string, 0, len(config.Order.Categories))
		for _, catName := range config.Order.Categories {
			cat, ok := config.Category[catName]
			if !ok {
				continue
			}
			categoryActive := resolver.Matches(cat.Tags)
			if activeOnly && !categoryActive {
				continue
			}
			visibleCategories = append(visibleCategories, catName)
		}

		for catIndex, catName := range visibleCategories {
			cat := config.Category[catName]
			categoryActive := resolver.Matches(cat.Tags)
			categoryLast := catIndex == len(visibleCategories)-1
			categoryPrefix := childPrefix("", categoryLast)
			fmt.Printf("%s%s %s\n", branch(categoryLast), marker(categoryActive), catName)

			visibleModules := make([]toml.Module, 0, len(cat.Modules))
			for _, mod := range cat.Modules {
				moduleActive := categoryActive && resolver.Matches(mod.Tags)
				if activeOnly && !moduleActive {
					continue
				}
				visibleModules = append(visibleModules, mod)
			}

			for modIndex, mod := range visibleModules {
				moduleActive := categoryActive && resolver.Matches(mod.Tags)
				moduleLast := modIndex == len(visibleModules)-1
				referencePath := module.NormalizeReferencePath(mod.Source)

				if showTags {
					var tags []string
					filePath := config.ResolveModulePath(mod.Source)
					fileTags := extractTagsFromFile(filePath)
					if len(fileTags) > 0 {
						tags = fileTags
					} else if len(mod.Tags) > 0 {
						tags = mod.Tags
					}

					fmt.Printf("%s%s%s %s", categoryPrefix, branch(moduleLast), marker(moduleActive), mod.Source)
					printTags(tags)
				} else {
					fmt.Printf("%s%s%s %s\n", categoryPrefix, branch(moduleLast), marker(moduleActive), mod.Source)
				}

				filePath := config.ResolveModulePath(mod.Source)
				parsedModule, err := module.Parse(filePath)
				if err != nil {
					return fmt.Errorf("failed to parse module %s: %w", mod.Source, err)
				}
				selectedIDs := selectedReferences[referencePath]
				if len(selectedIDs) == 0 && len(config.Activate.References) == 0 {
					selectedIDs = selectedModuleBlocks(mod)
				}

				blockLines := make([]string, 0, len(parsedModule.Blocks))
				for i, block := range parsedModule.Blocks {
					if block.Metadata.ID == "" {
						continue
					}
					included := moduleActive && blockIncluded(parsedModule.Blocks, i, selectedIDs)
					if activeOnly && !included {
						continue
					}
					description := block.Metadata.TLDR
					if description != "" {
						description = "  " + description
					}
					blockLines = append(blockLines, fmt.Sprintf("%s %-28s %s%s", marker(included), block.Metadata.ID, block.Heading, description))
				}

				modulePrefix := childPrefix(categoryPrefix, moduleLast)
				for blockIndex, line := range blockLines {
					blockLast := blockIndex == len(blockLines)-1
					fmt.Printf("%s%s%s\n", modulePrefix, branch(blockLast), line)
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
	listCmd.Flags().StringP("config", "c", "AGENTS.toml", "Path to AGENTS.toml config file")
}

func selectedReferencesByPath(config *toml.Config) (map[string]map[string]bool, error) {
	selected := make(map[string]map[string]bool)
	for _, rawReference := range config.Activate.References {
		reference, err := module.ParseReference(rawReference)
		if err != nil {
			return nil, err
		}
		if selected[reference.Path] == nil {
			selected[reference.Path] = make(map[string]bool)
		}
		selected[reference.Path][reference.ID] = true
	}
	return selected, nil
}

func selectedModuleBlocks(mod toml.Module) map[string]bool {
	selected := make(map[string]bool)
	for _, blockID := range mod.Blocks {
		selected[strings.TrimSpace(blockID)] = true
	}
	return selected
}

func blockIncluded(blocks []module.Block, index int, selectedIDs map[string]bool) bool {
	block := blocks[index]
	if selectedIDs[block.Metadata.ID] {
		return true
	}
	for i := index - 1; i >= 0; i-- {
		ancestor := blocks[i]
		if ancestor.Level >= block.Level {
			continue
		}
		if selectedIDs[ancestor.Metadata.ID] {
			return true
		}
		block = ancestor
	}
	return false
}

func marker(active bool) string {
	if active {
		return "+"
	}
	return "-"
}

func branch(last bool) string {
	if last {
		return "└── "
	}
	return "├── "
}

func childPrefix(prefix string, parentLast bool) string {
	if parentLast {
		return prefix + "    "
	}
	return prefix + "│   "
}

func printTags(tags []string) {
	if len(tags) > 0 {
		fmt.Printf(" {%s}\n", strings.Join(tags, ", "))
		return
	}
	fmt.Println()
}
