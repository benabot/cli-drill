package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/benabot/cli-drill/data"
)

func TestDirectoryUsesScanCatalogWhenConfigIsProvided(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--config", configPath, "directory", "--type", "alias"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !bytes.Contains([]byte(got), []byte("cgs\talias")) {
		t.Fatalf("expected scanned alias in directory output, got:\n%s", got)
	}
}

func writeAppTestConfig(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "config.toml")
	content := `dotfiles_path = "` + root + `"
shell = "zsh"

[paths]
aliases = ["zsh/modules/aliases.zsh"]
functions = []
docs = []

[security]
exclude = []
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile config returned error: %v", err)
	}
	return path
}

func writeAppTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
