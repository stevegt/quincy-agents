// Intent: Assembly engine that combines modules based on active scopes,
// category order, and stable heading-subtree selections. Source: DI-lorad

package assemble

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/scope"
	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
)

type Engine struct {
	config   *toml.Config
	resolver *scope.Resolver
	warnings []string
}

func NewEngine(config *toml.Config) *Engine {
	return &Engine{
		config:   config,
		resolver: scope.NewResolver(config.Activate.Scopes),
	}
}

func (e *Engine) Assemble() (string, error) {
	return e.assembleBlocks()
}

func (e *Engine) Warnings() []string {
	return append([]string(nil), e.warnings...)
}

func (e *Engine) assembleBlocks() (string, error) {
	sources := e.moduleSources()
	if len(sources) == 0 {
		return "", fmt.Errorf("no active module sources configured")
	}
	index, err := module.BuildIndex(sources)
	if err != nil {
		return "", fmt.Errorf("build module index: %w", err)
	}

	requested, err := e.requestedReferences()
	if err != nil {
		return "", err
	}
	if len(requested) == 0 {
		return "", fmt.Errorf("no block selections configured; add activate.references or module blocks")
	}

	var sections []string
	for _, reference := range requested {
		block, err := index.SelectBlock(reference)
		if err != nil {
			return "", err
		}
		// Intent: Preserve build momentum for heading-only blocks while warning
		// that the selected block has no body content. Source: DI-nasot
		if block.Block.Content != "" {
			sections = append(sections, block.Block.Content)
		}
		if blockBodyEmpty(block.Block.Content) {
			e.warnings = append(e.warnings, fmt.Sprintf("selected block %s has no body content", reference.String()))
		}
	}

	if strings.TrimSpace(strings.Join(sections, "\n\n")) == "" {
		return "", fmt.Errorf("rendered output is empty")
	}

	return strings.Join(sections, "\n\n"), nil
}

func ValidateOutput(output string) error {
	if strings.TrimSpace(output) == "" {
		return fmt.Errorf("rendered output is empty")
	}
	return nil
}

func blockBodyEmpty(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		return false
	}
	return true
}

func (e *Engine) moduleSources() []module.SourceFile {
	var sources []module.SourceFile
	seen := make(map[string]bool)
	for _, catName := range e.categoryOrder() {
		cat, ok := e.config.Category[catName]
		if !ok || !e.resolver.Matches(cat.Tags) {
			continue
		}
		for _, mod := range cat.Modules {
			if !e.resolver.Matches(mod.Tags) {
				continue
			}
			referencePath := module.NormalizeReferencePath(mod.Source)
			filePath := e.config.ResolveModulePath(mod.Source)
			key := filepath.Clean(filePath)
			if seen[key] {
				continue
			}
			seen[key] = true
			sources = append(sources, module.SourceFile{
				ReferencePath: referencePath,
				FilePath:      filePath,
			})
		}
	}
	return sources
}

func (e *Engine) requestedReferences() ([]module.Reference, error) {
	if len(e.config.Activate.References) > 0 {
		references := make([]module.Reference, 0, len(e.config.Activate.References))
		for _, rawReference := range e.config.Activate.References {
			reference, err := module.ParseReference(rawReference)
			if err != nil {
				return nil, err
			}
			references = append(references, reference)
		}
		return references, nil
	}

	var references []module.Reference
	for _, catName := range e.categoryOrder() {
		cat, ok := e.config.Category[catName]
		if !ok || !e.resolver.Matches(cat.Tags) {
			continue
		}
		for _, mod := range cat.Modules {
			if !e.resolver.Matches(mod.Tags) {
				continue
			}
			referencePath := module.NormalizeReferencePath(mod.Source)
			for _, blockID := range mod.Blocks {
				references = append(references, module.Reference{
					Path: referencePath,
					ID:   strings.TrimSpace(blockID),
				})
			}
		}
	}
	return references, nil
}

func (e *Engine) categoryOrder() []string {
	categoryOrder := e.config.Order.Categories
	if len(categoryOrder) == 0 {
		for catName := range e.config.Category {
			categoryOrder = append(categoryOrder, catName)
		}
	}
	return categoryOrder
}
