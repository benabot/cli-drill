package detect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benabot/cli-drill/internal/catalog"
	"github.com/benabot/cli-drill/internal/config"
)

func TestScanDeduplicatesToolDetectedFromAliasAndDocs(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "zsh/modules/aliases.zsh", "alias findmd='fd -e md | fzf'\n")
	writeFile(t, root, "docs/tools.md", "Use fd with fzf.\n")

	cfg := config.Default()
	cfg.DotfilesPath = root
	cfg.Paths = config.PathConfig{
		Aliases: []string{"zsh/modules/aliases.zsh"},
		Docs:    []string{"docs/tools.md"},
	}

	entries, warnings := Scan(cfg)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %#v", warnings)
	}

	tools := entries.FilterByType(catalog.EntryTool)
	countFD := 0
	for _, entry := range tools {
		if entry.ID == "tool-fd" {
			countFD++
			if len(entry.Sources) != 2 {
				t.Fatalf("expected merged fd sources, got %#v", entry.Sources)
			}
		}
	}
	if countFD != 1 {
		t.Fatalf("expected one fd tool, got %d in %#v", countFD, tools)
	}
}

func TestScanSkipsGenericDuplicateMarkdownHeadings(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/tools.md", "# Tools\n## rg\n## rg\n## Notes\n## Atuin history search\n")

	cfg := config.Default()
	cfg.DotfilesPath = root
	cfg.Paths = config.PathConfig{Docs: []string{"docs/tools.md"}}

	entries, warnings := Scan(cfg)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %#v", warnings)
	}

	concepts := entries.FilterByType(catalog.EntryConcept)
	if len(concepts) != 1 {
		t.Fatalf("expected one useful concept, got %d: %#v", len(concepts), concepts)
	}
	if concepts[0].Name != "Atuin history search" {
		t.Fatalf("unexpected concept: %#v", concepts[0])
	}
}

func TestResolveConfiguredPathRejectsAbsolutePathOutsideDotfilesRoot(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.md")

	if _, err := ResolveConfiguredPath(root, outside); err == nil {
		t.Fatal("expected absolute path outside dotfiles root to be rejected")
	}
}

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
