// Intent: Provide a navigable selector shell before growing the module library,
// so users can inspect block choices without hand-editing TOML first. Source: DI-bakom

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/scope"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open the module selector TUI",
	Long:  "Browse configured modules and heading blocks in a yazi-like terminal selector",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		if configFile == "" {
			configFile = "AGENTS.toml"
		}

		config, err := toml.Parse(configFile)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", configFile, err)
		}

		model := newTUIModel(config)
		if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("run tui: %w", err)
		}
		return nil
	},
}

func init() {
	tuiCmd.Flags().StringP("config", "c", "AGENTS.toml", "Path to AGENTS.toml config file")
}

type tuiItemKind int

const (
	tuiItemCategory tuiItemKind = iota
	tuiItemModule
	tuiItemBlock
	tuiItemDiagnostic
)

type tuiItem struct {
	kind         tuiItemKind
	category     string
	source       string
	blockID      string
	heading      string
	description  string
	indent       int
	state        selectionState
	toggleable   bool
	moduleActive bool
	categoryLast bool
	moduleLast   bool
	blockLast    bool
}

type tuiModel struct {
	config        *toml.Config
	items         []tuiItem
	selected      map[string]map[string]bool
	cursor        int
	width         int
	height        int
	statusMessage string
}

func newTUIModel(config *toml.Config) tuiModel {
	selected, err := selectedReferencesByPath(config)
	if err != nil {
		selected = map[string]map[string]bool{}
	}
	model := tuiModel{
		config:   config,
		selected: selected,
		width:    100,
		height:   30,
	}
	model.rebuildItems()
	return model
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "enter":
			m.toggleCurrentItem()
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	var out strings.Builder
	out.WriteString(headerStyle.Render("mogent module selector"))
	out.WriteString("\n")
	out.WriteString(helpStyle.Render("j/k move  space toggle block  q quit  writes are disabled in this POC"))
	out.WriteString("\n\n")

	listHeight := m.height - 5
	if listHeight < 1 {
		listHeight = len(m.items)
	}
	start, end := visibleWindow(m.cursor, len(m.items), listHeight)
	for i := start; i < end; i++ {
		line := m.renderItem(i)
		out.WriteString(truncateToWidth(line, m.width))
		out.WriteString("\n")
	}

	if m.statusMessage != "" {
		out.WriteString("\n")
		out.WriteString(statusStyle.Render(m.statusMessage))
	}
	return out.String()
}

func (m *tuiModel) rebuildItems() {
	resolver := scope.NewResolver(m.config.Activate.Scopes)
	items := make([]tuiItem, 0)
	categoryOrder := m.config.Order.Categories
	if len(categoryOrder) == 0 {
		for category := range m.config.Category {
			categoryOrder = append(categoryOrder, category)
		}
	}

	for categoryIndex, categoryName := range categoryOrder {
		category, ok := m.config.Category[categoryName]
		if !ok {
			items = append(items, tuiItem{
				kind:        tuiItemDiagnostic,
				heading:     fmt.Sprintf("category %q is in order but not defined", categoryName),
				description: "AGENTS.toml",
			})
			continue
		}

		categoryActive := resolver.Matches(category.Tags)
		categoryLast := categoryIndex == len(categoryOrder)-1
		items = append(items, tuiItem{
			kind:         tuiItemCategory,
			category:     categoryName,
			heading:      categoryName,
			indent:       0,
			state:        boolSelectionState(categoryActive),
			categoryLast: categoryLast,
		})

		for moduleIndex, configuredModule := range category.Modules {
			moduleActive := categoryActive && resolver.Matches(configuredModule.Tags)
			moduleLast := moduleIndex == len(category.Modules)-1
			referencePath := module.NormalizeReferencePath(configuredModule.Source)
			items = append(items, tuiItem{
				kind:         tuiItemModule,
				category:     categoryName,
				source:       configuredModule.Source,
				heading:      configuredModule.Source,
				description:  configuredModule.Name,
				indent:       1,
				state:        boolSelectionState(moduleActive),
				moduleActive: moduleActive,
				categoryLast: categoryLast,
				moduleLast:   moduleLast,
			})

			selectedIDs := m.selected[referencePath]
			if len(selectedIDs) == 0 && len(m.config.Activate.References) == 0 {
				selectedIDs = selectedModuleBlocks(configuredModule)
				if m.selected[referencePath] == nil {
					m.selected[referencePath] = selectedIDs
				}
			}

			modulePath := m.config.ResolveModulePath(configuredModule.Source)
			parsedModule, err := module.Parse(modulePath)
			if err != nil {
				items = append(items, tuiItem{
					kind:         tuiItemDiagnostic,
					category:     categoryName,
					source:       configuredModule.Source,
					heading:      fmt.Sprintf("cannot parse %s", configuredModule.Source),
					description:  err.Error(),
					indent:       2,
					state:        selectionInactive,
					categoryLast: categoryLast,
					moduleLast:   moduleLast,
				})
				if !filepath.IsAbs(configuredModule.Source) && filepath.Ext(configuredModule.Source) == "" {
					items = append(items, tuiItem{
						kind:         tuiItemDiagnostic,
						category:     categoryName,
						source:       configuredModule.Source,
						heading:      "module source has no extension",
						description:  "did you mean " + configuredModule.Source + ".md?",
						indent:       2,
						state:        selectionInactive,
						categoryLast: categoryLast,
						moduleLast:   moduleLast,
					})
				}
				continue
			}

			for blockIndex, block := range parsedModule.Blocks {
				if block.Metadata.ID == "" {
					continue
				}
				state := blockSelectionState(parsedModule.Blocks, blockIndex, selectedIDs, moduleActive)
				items = append(items, tuiItem{
					kind:         tuiItemBlock,
					category:     categoryName,
					source:       referencePath,
					blockID:      block.Metadata.ID,
					heading:      block.Heading,
					description:  block.Metadata.TLDR,
					indent:       2 + block.Level - 1,
					state:        state,
					toggleable:   moduleActive,
					moduleActive: moduleActive,
					categoryLast: categoryLast,
					moduleLast:   moduleLast,
					blockLast:    blockIndex == len(parsedModule.Blocks)-1,
				})
			}

			if !filepath.IsAbs(configuredModule.Source) && filepath.Ext(configuredModule.Source) == "" {
				items = append(items, tuiItem{
					kind:         tuiItemDiagnostic,
					category:     categoryName,
					source:       configuredModule.Source,
					heading:      "module source has no extension",
					description:  "did you mean " + configuredModule.Source + ".md?",
					indent:       2,
					state:        selectionInactive,
					categoryLast: categoryLast,
					moduleLast:   moduleLast,
				})
			}
		}
	}

	m.items = items
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *tuiModel) toggleCurrentItem() {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return
	}
	item := m.items[m.cursor]
	if item.kind != tuiItemBlock || !item.toggleable {
		m.statusMessage = "Only active blocks can be toggled in this POC."
		return
	}
	if m.selected[item.source] == nil {
		m.selected[item.source] = make(map[string]bool)
	}
	if m.selected[item.source][item.blockID] {
		delete(m.selected[item.source], item.blockID)
		m.statusMessage = fmt.Sprintf("Removed [[%s#%s]] from the in-memory selection.", item.source, item.blockID)
	} else {
		m.selected[item.source][item.blockID] = true
		m.statusMessage = fmt.Sprintf("Added [[%s#%s]] to the in-memory selection.", item.source, item.blockID)
	}
	m.rebuildItems()
}

func (m tuiModel) renderItem(index int) string {
	item := m.items[index]
	cursor := " "
	if index == m.cursor {
		cursor = ">"
	}

	line := fmt.Sprintf("%s %s %s%s", cursor, selectionMarker(item.state), strings.Repeat("  ", item.indent), item.heading)
	if item.kind == tuiItemBlock && item.blockID != "" {
		line = fmt.Sprintf("%s %s %s%-24s %s", cursor, selectionMarker(item.state), strings.Repeat("  ", item.indent), item.blockID, item.heading)
	}
	if item.description != "" && m.width >= 80 {
		line += "  " + item.description
	}

	switch item.kind {
	case tuiItemCategory:
		return lipgloss.NewStyle().Bold(true).Render(line)
	case tuiItemModule:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render(line)
	case tuiItemDiagnostic:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render(line)
	default:
		if index == m.cursor {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("237")).Render(line)
		}
		return line
	}
}

func boolSelectionState(active bool) selectionState {
	if active {
		return selectionExplicit
	}
	return selectionInactive
}

func visibleWindow(cursor int, total int, height int) (int, int) {
	if total <= height {
		return 0, total
	}
	half := height / 2
	start := cursor - half
	if start < 0 {
		start = 0
	}
	if start+height > total {
		start = total - height
	}
	return start, start + height
}
