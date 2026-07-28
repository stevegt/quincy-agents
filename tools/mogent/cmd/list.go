// Intent: Render module selections as a tree so active/inactive blocks are
// readable before the project grows a richer TUI.
// Source: DI-soviv

package cmd

import (
	"fmt"
	"os"
	"strconv"
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
		noDescriptions, _ := cmd.Flags().GetBool("no-descriptions")
		widthOverride, _ := cmd.Flags().GetInt("width")
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
		// Intent: Keep tree output readable in narrow terminals while retaining
		// explicit, inherited, and inactive selection states. Source: DI-nasot
		width := terminalWidth(widthOverride)
		showDescriptions := !noDescriptions && width >= 72

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
					state := blockSelectionState(parsedModule.Blocks, i, selectedIDs, moduleActive)
					if activeOnly && state == selectionInactive {
						continue
					}
					description := block.Metadata.TLDR
					if showDescriptions && description != "" {
						description = "  " + description
					} else {
						description = ""
					}
					blockLines = append(blockLines, fmt.Sprintf("%s %-28s %s%s", selectionMarker(state), block.Metadata.ID, block.Heading, description))
				}

				modulePrefix := childPrefix(categoryPrefix, moduleLast)
				for blockIndex, line := range blockLines {
					blockLast := blockIndex == len(blockLines)-1
					prefix := modulePrefix + branch(blockLast)
					fmt.Println(prefix + truncateToWidth(line, width-displayWidth(prefix)))
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
	listCmd.Flags().Bool("no-descriptions", false, "Hide block descriptions")
	listCmd.Flags().Int("width", 0, "Render width override for testing")
	listCmd.Flags().StringP("config", "c", "AGENTS.toml", "Path to AGENTS.toml config file")
}

type selectionState int

const (
	selectionInactive selectionState = iota
	selectionInherited
	selectionExplicit
)

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

func blockSelectionState(blocks []module.Block, index int, selectedIDs map[string]bool, moduleActive bool) selectionState {
	if !moduleActive {
		return selectionInactive
	}
	block := blocks[index]
	if selectedIDs[block.Metadata.ID] {
		return selectionExplicit
	}
	for i := index - 1; i >= 0; i-- {
		ancestor := blocks[i]
		if ancestor.Level >= block.Level {
			continue
		}
		if selectedIDs[ancestor.Metadata.ID] {
			return selectionInherited
		}
		block = ancestor
	}
	return selectionInactive
}

func marker(active bool) string {
	if active {
		return "+"
	}
	return "-"
}

func selectionMarker(state selectionState) string {
	switch state {
	case selectionExplicit:
		return "+"
	case selectionInherited:
		return "|"
	default:
		return "-"
	}
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

func terminalWidth(override int) int {
	if override > 0 {
		return override
	}
	columns := strings.TrimSpace(os.Getenv("COLUMNS"))
	if columns != "" {
		width, err := strconv.Atoi(columns)
		if err == nil && width > 0 {
			return width
		}
	}
	return 100
}

func truncateToWidth(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if displayWidth(value) <= width {
		return value
	}
	if width <= 3 {
		return strings.Repeat(".", width)
	}
	runes := []rune(value)
	for displayWidth(string(runes)) > width-3 {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "..."
}

func displayWidth(value string) int {
	return len([]rune(value))
}
