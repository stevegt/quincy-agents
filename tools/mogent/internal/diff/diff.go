// Intent: Show which selectable heading blocks are included/excluded by active
// scopes so block-aware diffs can report stable IDs and rendered text impact.
// Source: DI-lorad

package diff

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/scope"
)

type SectionStatus struct {
	Name      string
	BlockID   string
	Included  bool
	Tags      []string
	LineCount int
}

type DiffResult struct {
	Sections []SectionStatus
	Stats    Stats
}

type Stats struct {
	Included int
	Excluded int
}

func AnalyzeModule(mod *module.Module, scopes []string) *DiffResult {
	resolver := scope.NewResolver(scopes)
	var sections []SectionStatus
	included := 0
	excluded := 0

	for _, block := range mod.Blocks {
		if block.Type != "header" {
			continue
		}

		status := SectionStatus{
			Name:      block.Heading,
			BlockID:   block.Metadata.ID,
			Tags:      block.Tags,
			LineCount: renderedLineCount(block.Content),
		}

		if len(block.Tags) == 0 {
			status.Included = true
			included++
		} else if resolver.Matches(block.Tags) {
			status.Included = true
			included++
		} else {
			status.Included = false
			excluded++
		}

		sections = append(sections, status)
	}

	return &DiffResult{
		Sections: sections,
		Stats: Stats{
			Included: included,
			Excluded: excluded,
		},
	}
}

func FormatAnalysis(result *DiffResult, moduleName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Module: %s\n", moduleName))
	sb.WriteString(fmt.Sprintf("Included: %d, Excluded: %d\n\n", result.Stats.Included, result.Stats.Excluded))

	for _, section := range result.Sections {
		blockSuffix := ""
		if section.BlockID != "" {
			blockSuffix = fmt.Sprintf(" #%s", section.BlockID)
		}
		if section.Included {
			sb.WriteString(fmt.Sprintf("  [+] %s%s", section.Name, blockSuffix))
			if len(section.Tags) > 0 {
				sb.WriteString(fmt.Sprintf(" {%s}", strings.Join(section.Tags, " ")))
			}
			sb.WriteString(fmt.Sprintf(" (%d rendered lines)\n", section.LineCount))
		} else {
			sb.WriteString(fmt.Sprintf("  [-] %s%s {%s} (%d rendered lines)\n", section.Name, blockSuffix, strings.Join(section.Tags, " "), section.LineCount))
		}
	}

	return sb.String()
}

func renderedLineCount(content string) int {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}
