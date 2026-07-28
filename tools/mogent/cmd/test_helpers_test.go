package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFileForTest(t *testing.T, dir string, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
