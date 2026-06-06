package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

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
	if !bytes.Contains([]byte(got), []byte("alias\tcgs\tcgs")) {
		t.Fatalf("expected scanned alias in directory output, got:\n%s", got)
	}
}

func TestScanSummaryPrintsCountsWithoutEntries(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	out, err := executeTestCommand("--config", configPath, "scan", "--summary")
	if err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out)
	}

	for _, want := range []string{"aliases: 1", "functions: 0", "tools: 0", "concepts: 0", "workflows: 0", "chapters: 0", "total: 1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected summary to contain %q, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "cgs\talias") {
		t.Fatalf("summary should not print entries, got:\n%s", out)
	}
}

func TestScanTypeFiltersAndRejectsUnknownType(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	out, err := executeTestCommand("--config", configPath, "scan", "--type", "alias")
	if err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "alias\tcgs\tcgs\tgit status --short") {
		t.Fatalf("expected alias output, got:\n%s", out)
	}

	out, err = executeTestCommand("--config", configPath, "scan", "--type", "nope")
	if err == nil {
		t.Fatalf("expected invalid type error, got output:\n%s", out)
	}
	if !strings.Contains(err.Error(), "unknown entry type") {
		t.Fatalf("expected clear invalid type error, got: %v", err)
	}
}

func TestShowFindsEntryByIDAndName(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	for _, query := range []string{"cgs", "cgs"} {
		out, err := executeTestCommand("--config", configPath, "show", query)
		if err != nil {
			t.Fatalf("Execute returned error: %v\n%s", err, out)
		}
		for _, want := range []string{"id: cgs", "name: cgs", "type: alias", "command: git status --short", "source: zsh/modules/aliases.zsh:1"} {
			if !strings.Contains(out, want) {
				t.Fatalf("expected show output to contain %q, got:\n%s", want, out)
			}
		}
	}
}

func TestShowReportsAmbiguousName(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias dup='git status --short'\n")
	writeAppTestFile(t, root, "zsh/modules/functions.zsh", "dup() {\n  echo dup\n}\n")
	configPath := writeAppTestConfigWithFunctions(t, root)

	out, err := executeTestCommand("--config", configPath, "show", "dup")
	if err == nil {
		t.Fatalf("expected ambiguous show error, got output:\n%s", out)
	}
	if !strings.Contains(err.Error(), "ambiguous entry") {
		t.Fatalf("expected ambiguous entry error, got: %v", err)
	}
}

func TestTrainMatchesTerminalControlShortcutInput(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl A?
    answer:
      primary: Ctrl+A
  - id: ctrl-e
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl E?
    answer:
      primary: Ctrl+E
  - id: ctrl-l
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl L?
    answer:
      primary: Ctrl+L
  - id: ctrl-l-bare
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl L bare?
    answer:
      primary: Ctrl+L
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("\x01\n\x05\n\x0c\nL\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Count(got, "Correct.") != 3 {
		t.Fatalf("expected three control shortcut answers to be correct, got:\n%s", got)
	}
	if !strings.Contains(got, "Pas encore. Reponse attendue: Ctrl+L") {
		t.Fatalf("expected bare L to be rejected, got:\n%s", got)
	}
}

func writeAppTestConfigWithFunctions(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "config.toml")
	content := `dotfiles_path = "` + root + `"
shell = "zsh"

[paths]
aliases = ["zsh/modules/aliases.zsh"]
functions = ["zsh/modules/functions.zsh"]
docs = []

[security]
exclude = []
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile config returned error: %v", err)
	}
	return path
}

func TestInitPrintWritesDefaultConfigToStdoutOnly(t *testing.T) {
	out, err := executeTestCommand("init", "--print")
	if err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out)
	}
	if !strings.Contains(out, `dotfiles_path = "~/dotfiles"`) || !strings.Contains(out, `shell = "zsh"`) {
		t.Fatalf("expected default TOML config, got:\n%s", out)
	}
	if strings.Contains(out, "Created config:") {
		t.Fatalf("init --print should not create a file, got:\n%s", out)
	}
}

func executeTestCommand(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
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
