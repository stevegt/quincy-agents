package module

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseHeadingSubtreeWithAgentModuleMetadata(t *testing.T) {
	path := writeTestModule(t, `# Thinking Style

<!--
agent_module:
  id: thinking-style
  tldr: Choose reasoning depth.
-->
Intro.

## Think Light

<!--
agent_module:
  id: think-light
  tldr: Prefer fast practical answers.
-->
Use fast practical reasoning.

### Socratic Mode

<!--
agent_module:
  id: socratic
-->
Ask short guiding questions.

## Think Hard

<!--
agent_module:
  id: think-hard
-->
Compare alternatives.
`)

	mod, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(mod.Blocks) != 4 {
		t.Fatalf("len(Blocks) = %d, want 4", len(mod.Blocks))
	}

	thinkLight := mod.Blocks[1]
	if thinkLight.Metadata.ID != "think-light" {
		t.Fatalf("Metadata.ID = %q, want think-light", thinkLight.Metadata.ID)
	}
	if !strings.Contains(thinkLight.Content, "### Socratic Mode") {
		t.Fatal("Think Light subtree did not include descendant heading")
	}
	if strings.Contains(thinkLight.Content, "agent_module:") {
		t.Fatal("Think Light rendered content retained builder metadata")
	}
	if strings.Contains(thinkLight.Content, "## Think Hard") {
		t.Fatal("Think Light subtree included sibling heading")
	}
}

func TestBuildIndexRejectsMissingMetadataID(t *testing.T) {
	path := writeTestModule(t, `# Missing ID

No metadata here.
`)

	_, err := BuildIndex([]SourceFile{{
		ReferencePath: "missing",
		FilePath:      path,
	}})
	if err == nil {
		t.Fatal("BuildIndex returned nil error for missing metadata ID")
	}
}

func TestBuildIndexRejectsDuplicateReferences(t *testing.T) {
	path := writeTestModule(t, `# First

<!--
agent_module:
  id: duplicate
-->
First.

# Second

<!--
agent_module:
  id: duplicate
-->
Second.
`)

	_, err := BuildIndex([]SourceFile{{
		ReferencePath: "duplicates",
		FilePath:      path,
	}})
	if err == nil {
		t.Fatal("BuildIndex returned nil error for duplicate references")
	}
}

func writeTestModule(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "module.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	return path
}
