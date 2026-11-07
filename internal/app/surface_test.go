package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/xdg"
)

func TestHelpExposesOnlyTrainingCommands(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{"chapters": true, "train": true, "stats": true, "reset": true, "completion": true, "help": true}
	for _, child := range cmd.Commands() {
		if !allowed[child.Name()] {
			t.Errorf("obsolete command is still public: %s", child.Name())
		}
	}
	for _, term := range []string{"--config", "dotfiles", "catalog", "scan"} {
		if strings.Contains(strings.ToLower(out.String()), term) {
			t.Errorf("obsolete help reference %q: %s", term, out.String())
		}
	}
	for _, name := range []string{"chapters", "train", "stats", "reset", "--progress"} {
		if !strings.Contains(out.String(), name) {
			t.Errorf("missing training interface %q", name)
		}
	}
}

func TestRemovedInterfaceRejectsRequestsAndPreservesLegacyFiles(t *testing.T) {
	isolateChapterHome(t)
	files := map[string]string{
		filepath.Join(legacyConfigRoot(t), "config.toml"):          "[broken legacy TOML",
		filepath.Join(legacyConfigRoot(t), "chapters", "old.yaml"): "items: [",
		xdg.ProgressFile(): `{"version":1,"completed_exercises":["old/item"],"chapter_scores":{},"last_session":null,"streak":0}`,
	}
	for path, content := range files {
		writeAppTestFile(t, filepath.Dir(path), filepath.Base(path), content)
	}
	for _, name := range []string{"init", "scan", "generate", "directory", "search", "show", "--config"} {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := NewRootCommand(data.Chapters())
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			cmd.SetArgs([]string{name})
			err := cmd.Execute()
			want := "unknown command"
			if name == "--config" {
				want = "unknown flag"
			}
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s was not rejected as removed: %v", name, err)
			}
		})
	}
	for path, original := range files {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != original {
			t.Errorf("legacy file changed: %s (%v)", path, err)
		}
	}
}
