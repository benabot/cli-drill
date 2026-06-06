package catalog

import "testing"

func TestAddDeduplicatesByTypeAndID(t *testing.T) {
	c := New()
	c.Add(Entry{
		ID:      "tool-rg",
		Name:    "rg",
		Type:    EntryTool,
		Summary: "Detected from alias command",
		Source:  Source{Path: "zsh/aliases.zsh", Line: 2},
		Tags:    []string{"alias"},
	})
	c.Add(Entry{
		ID:      "tool-rg",
		Name:    "rg",
		Type:    EntryTool,
		Summary: "Detected in documentation",
		Source:  Source{Path: "docs/tools.md", Line: 8},
		Tags:    []string{"docs", "alias"},
	})

	entries := c.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 deduplicated entry, got %d: %#v", len(entries), entries)
	}

	got := entries[0]
	if got.Summary != "Detected from alias command; Detected in documentation" {
		t.Fatalf("unexpected merged summary: %q", got.Summary)
	}
	if len(got.Sources) != 2 {
		t.Fatalf("expected 2 merged sources, got %#v", got.Sources)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "alias" || got.Tags[1] != "docs" {
		t.Fatalf("unexpected merged tags: %#v", got.Tags)
	}
}

func TestAddKeepsDifferentTypesWithSameID(t *testing.T) {
	c := New(
		Entry{ID: "rg", Name: "rg", Type: EntryConcept},
		Entry{ID: "rg", Name: "rg", Type: EntryTool},
	)

	if got := len(c.Entries()); got != 2 {
		t.Fatalf("expected different types to remain distinct, got %d", got)
	}
}
