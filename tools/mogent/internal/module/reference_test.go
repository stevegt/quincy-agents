package module

import "testing"

func TestParseReference(t *testing.T) {
	reference, err := ParseReference("[[cognition/thinking-style.md#think-light]]")
	if err != nil {
		t.Fatalf("ParseReference returned error: %v", err)
	}

	if reference.Path != "cognition/thinking-style" {
		t.Fatalf("Path = %q, want cognition/thinking-style", reference.Path)
	}
	if reference.ID != "think-light" {
		t.Fatalf("ID = %q, want think-light", reference.ID)
	}
}

func TestParseReferenceRejectsMissingBlockID(t *testing.T) {
	if _, err := ParseReference("[[cognition/thinking-style]]"); err == nil {
		t.Fatal("ParseReference returned nil error for missing block ID")
	}
}
