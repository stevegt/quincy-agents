package cmd

import (
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
)

func TestTruncateToWidth(t *testing.T) {
	got := truncateToWidth("abcdefghijklmnopqrstuvwxyz", 10)
	if got != "abcdefg..." {
		t.Fatalf("truncateToWidth = %q, want abcdefg...", got)
	}
}

func TestTruncateToWidthHidesWhenTooSmall(t *testing.T) {
	got := truncateToWidth("abcdef", 2)
	if got != ".." {
		t.Fatalf("truncateToWidth = %q, want ..", got)
	}
}

func TestBlockSelectionStateDistinguishesExplicitAndInherited(t *testing.T) {
	blocks := []module.Block{
		{Level: 1, Metadata: module.Metadata{ID: "parent"}},
		{Level: 2, Metadata: module.Metadata{ID: "child"}},
		{Level: 1, Metadata: module.Metadata{ID: "sibling"}},
	}
	selectedIDs := map[string]bool{"parent": true}

	if got := blockSelectionState(blocks, 0, selectedIDs, true); got != selectionExplicit {
		t.Fatalf("parent state = %v, want explicit", got)
	}
	if got := blockSelectionState(blocks, 1, selectedIDs, true); got != selectionInherited {
		t.Fatalf("child state = %v, want inherited", got)
	}
	if got := blockSelectionState(blocks, 2, selectedIDs, true); got != selectionInactive {
		t.Fatalf("sibling state = %v, want inactive", got)
	}
}

func TestDiscoverModuleSources(t *testing.T) {
	dir := t.TempDir()
	writeFileForTest(t, dir, "identity.md")
	writeFileForTest(t, dir, "instructions.md")
	writeFileForTest(t, dir, "notes.txt")

	got := strings.Join(discoverModuleSources(dir), ",")
	if got != "identity.md,instructions.md" {
		t.Fatalf("discoverModuleSources = %q, want identity.md,instructions.md", got)
	}
}
