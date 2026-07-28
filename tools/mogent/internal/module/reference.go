// Intent: Parse stable Obsidian-style block references used by presets and
// config selections. Source: DI-lorad

package module

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Reference struct {
	Path string
	ID   string
}

func ParseReference(reference string) (Reference, error) {
	trimmed := strings.TrimSpace(reference)
	if strings.HasPrefix(trimmed, "[[") && strings.HasSuffix(trimmed, "]]") {
		trimmed = strings.TrimSuffix(strings.TrimPrefix(trimmed, "[["), "]]")
	}

	path, blockID, ok := strings.Cut(trimmed, "#")
	if !ok || strings.TrimSpace(path) == "" || strings.TrimSpace(blockID) == "" {
		return Reference{}, fmt.Errorf("reference %q must use [[path#id]] syntax", reference)
	}

	return Reference{
		Path: NormalizeReferencePath(path),
		ID:   strings.TrimSpace(blockID),
	}, nil
}

func NormalizeReferencePath(path string) string {
	cleaned := filepath.ToSlash(filepath.Clean(strings.TrimSpace(path)))
	return strings.TrimSuffix(cleaned, ".md")
}

func (r Reference) String() string {
	return fmt.Sprintf("[[%s#%s]]", r.Path, r.ID)
}
