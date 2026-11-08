package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"gopkg.in/yaml.v3"
)

func TestCLICorrectionWithEmbeddedExercises(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		chapter, item, input string
		correct              bool
	}{
		{"02-navigation-shell", "fzf-pick", " FZF ", true},
		{"02-navigation-shell", "fzf-pick", "\tfzf\t", true},
		{"02-navigation-shell", "zoxide-jump", "z", true},
		{"02-navigation-shell", "fd-find-dir", "fd", true},
		{"02-navigation-shell", "fd-find-dir", "bat", false},
		{"02-navigation-shell", "fd-fzf-dir", "fd --type d | fzf", true},
		{"02-navigation-shell", "fd-fzf-dir", "fd   -t d|fzf", true},
		{"02-navigation-shell", "fd-fzf-dir", "fd -t D | fzf", false},
		{"07-lecture-preview", "bat-readme", "bat readme.md", false},
		{"07-lecture-preview", "bat-readme", " bat README.md ", true},
		{"08-micro", "micro-save", "^s", true},
		{"08-micro", "micro-save", "s", false},
		{"05-outils-quotidiens", "atuin", "control-r", true},
		{"06-recherche-fichiers-contenu", "rg-fzf-todo", "\trg \"TODO\" | fzf\t", true},
		{"01-raccourcis-terminal", "ctrl-a", "\x01", true},
	} {
		t.Run(tc.chapter+"/"+tc.item+"/"+tc.input, func(t *testing.T) {
			var selected chapter.Chapter
			for _, ch := range chapters {
				for _, item := range ch.Items {
					if ch.ID == tc.chapter && item.ID == tc.item {
						selected = ch
						selected.Items = []chapter.Item{item}
					}
				}
			}
			if selected.ID == "" {
				t.Fatal("embedded exercise not found")
			}
			content, err := yaml.Marshal(selected)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(t.TempDir(), "progress.json")
			if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
				t.Fatal(err)
			}
			cmd := NewRootCommand(fstest.MapFS{"chapter.yaml": {Data: content}})
			var output bytes.Buffer
			cmd.SetIn(strings.NewReader(tc.input + "\n"))
			cmd.SetOut(&output)
			cmd.SetErr(&output)
			cmd.SetArgs([]string{"--progress", path, "train", tc.chapter})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			correct := 0
			var acquired []string
			if tc.correct {
				correct, acquired = 1, []string{tc.chapter + "/" + tc.item}
			}
			assertRecordingScore(t, path, tc.chapter, 1, correct, acquired...)
			if strings.Contains(output.String(), "Correct.") != tc.correct {
				t.Fatalf("feedback does not match persisted result: %s", output.String())
			}
		})
	}
}
