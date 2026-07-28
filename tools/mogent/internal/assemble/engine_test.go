package assemble

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
)

func TestAssembleSelectedBlocks(t *testing.T) {
	dir := t.TempDir()
	moduleDir := filepath.Join(dir, "modules")
	if err := os.Mkdir(moduleDir, 0755); err != nil {
		t.Fatalf("Mkdir returned error: %v", err)
	}
	modulePath := filepath.Join(moduleDir, "thinking.md")
	if err := os.WriteFile(modulePath, []byte(`# Thinking

<!--
agent_module:
  id: thinking
-->
Intro.

## Think Light

<!--
agent_module:
  id: think-light
-->
Use fast practical reasoning.

## Think Hard

<!--
agent_module:
  id: think-hard
-->
Compare alternatives.
`), 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	config := &toml.Config{
		Config: toml.ConfigSection{ModuleDir: moduleDir},
		Order:  toml.OrderSection{Categories: []string{"cognition"}},
		Category: map[string]toml.Category{
			"cognition": {
				Modules: []toml.Module{{
					Name:   "thinking",
					Source: "thinking.md",
					Blocks: []string{"think-hard"},
				}},
			},
		},
	}

	output, err := NewEngine(config).Assemble()
	if err != nil {
		t.Fatalf("Assemble returned error: %v", err)
	}
	if !strings.Contains(output, "## Think Hard") {
		t.Fatal("output did not include selected block")
	}
	if strings.Contains(output, "## Think Light") {
		t.Fatal("output included unselected sibling block")
	}
	if strings.Contains(output, "agent_module:") {
		t.Fatal("output retained builder metadata")
	}
}
