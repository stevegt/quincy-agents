package cmd

import (
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/toml"
)

func TestVisibleWindowCentersCursor(t *testing.T) {
	start, end := visibleWindow(10, 30, 5)
	if start != 8 || end != 13 {
		t.Fatalf("visibleWindow = %d,%d; want 8,13", start, end)
	}
}

func TestTUIModelReportsMissingModule(t *testing.T) {
	config := &toml.Config{
		Config: toml.ConfigSection{ModuleDir: t.TempDir()},
		Order:  toml.OrderSection{Categories: []string{"identity"}},
		Category: map[string]toml.Category{
			"identity": {
				Modules: []toml.Module{{Name: "identity", Source: "identity"}},
			},
		},
	}

	model := newTUIModel(config)
	model.width = 200
	got := model.View()
	if !strings.Contains(got, "cannot parse identity") {
		t.Fatalf("TUI view did not report missing module:\n%s", got)
	}
	if !strings.Contains(got, "did you mean identity.md?") {
		t.Fatalf("TUI view did not report extension hint:\n%s", got)
	}
}
