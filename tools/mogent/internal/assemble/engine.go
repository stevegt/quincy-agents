// Intent: Assembly engine that combines modules based on active scopes and category order.
// Source: DI-jusuk

package assemble

import (
	"os"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/scope"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
)

type Engine struct {
	config   *toml.Config
	resolver *scope.Resolver
}

func NewEngine(config *toml.Config) *Engine {
	return &Engine{
		config:   config,
		resolver: scope.NewResolver(config.Activate.Scopes),
	}
}

func (e *Engine) Assemble() (string, error) {
	var sections []string

	categoryOrder := e.config.Order.Categories
	if len(categoryOrder) == 0 {
		for catName := range e.config.Category {
			categoryOrder = append(categoryOrder, catName)
		}
	}

	for _, catName := range categoryOrder {
		cat, ok := e.config.Category[catName]
		if !ok {
			continue
		}

		if !e.resolver.Matches(cat.Tags) {
			continue
		}

		for _, mod := range cat.Modules {
			if !e.resolver.Matches(mod.Tags) {
				continue
			}

			path := e.config.ResolveModulePath(mod.Source)
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			filtered := e.filterContent(string(content), mod)
			if filtered != "" {
				sections = append(sections, filtered)
			}
		}
	}

	return strings.Join(sections, "\n\n"), nil
}

func (e *Engine) filterContent(content string, mod toml.Module) string {
	lines := strings.Split(content, "\n")
	var result []string
	inExcludedBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			tags := extractTagsFromLine(line)
			if len(tags) > 0 && !e.resolver.Matches(tags) {
				inExcludedBlock = true
				continue
			}
			inExcludedBlock = false
			line = stripTags(line)
		}

		if !inExcludedBlock {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func stripTags(line string) string {
	start := strings.Index(line, "{")
	end := strings.Index(line, "}")
	if start == -1 || end == -1 || end <= start {
		return line
	}
	return strings.TrimSpace(line[:start] + line[end+1:])
}

func extractTagsFromLine(line string) []string {
	start := strings.Index(line, "{")
	end := strings.Index(line, "}")
	if start == -1 || end == -1 || end <= start {
		return nil
	}

	tagStr := line[start+1 : end]
	var tags []string
	for _, t := range strings.Fields(tagStr) {
		t = strings.TrimPrefix(t, "#")
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}
