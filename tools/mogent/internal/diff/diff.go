// Intent: Show which sections of a module are included/excluded by active scopes.
// Source: DI-jusuk

package diff

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/scope"
)

type SectionStatus struct {
	Name    string
	Included bool
	Tags    []string
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
			Name: block.Content,
			Tags: block.Tags,
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
		if section.Included {
			sb.WriteString(fmt.Sprintf("  [+] %s", section.Name))
			if len(section.Tags) > 0 {
				sb.WriteString(fmt.Sprintf(" {%s}", strings.Join(section.Tags, " ")))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString(fmt.Sprintf("  [-] %s {%s}\n", section.Name, strings.Join(section.Tags, " ")))
		}
	}

	return sb.String()
}
