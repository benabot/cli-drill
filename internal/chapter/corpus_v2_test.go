package chapter

import (
	"io/fs"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/exercise"
	"gopkg.in/yaml.v3"
)

func TestV2Corpus(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(chapters) != 13 {
		t.Fatalf("want 13 chapters, got %d", len(chapters))
	}
	first := chapters[0]
	if first.ID != "00-premiers-pas" || first.Family != FamilyUnixShell || first.LearningOrder != LearningOrderGuided || len(first.Items) != 30 {
		t.Fatalf("unexpected first chapter: %#v", first)
	}
	answers := []string{"pwd", "ls", "ls -a", "cd docs", "cd ..", "cd /atelier", "mkdir essais", `cat "notes de cours.txt"`}
	for i, item := range first.Items[:len(answers)] {
		if item.ExerciseType != exercise.TypeSimpleShellSim || item.Difficulty != DifficultyBeginner || len(item.Tools) != 0 || item.Answer.Primary != answers[i] {
			t.Fatalf("unexpected first lesson %d: %#v", i, item)
		}
	}
	families := map[Family]int{}
	total := 0
	for _, ch := range chapters {
		if ch.Family == "" || ch.LearningOrder == "" {
			t.Fatalf("missing chapter metadata: %s", ch.ID)
		}
		for _, item := range ch.Items {
			// The existing TUI clips single lines at 80 columns. Leave room
			// for padding and the feedback prefix without changing its renderer.
			for _, text := range []struct {
				value string
				width int
			}{{item.Prompt, 70}, {item.Explanation, 60}} {
				for _, line := range strings.Split(text.value, "\n") {
					if utf8.RuneCountInString(line) > text.width {
						t.Fatalf("unreadable line: %s/%s: %q", ch.ID, item.ID, line)
					}
				}
			}
			if item.Difficulty == "" || item.Explanation == "" {
				t.Fatalf("incomplete exercise: %s/%s", ch.ID, item.ID)
			}
			if ch.Family == FamilyModernCLIMacOS && len(item.Tools) == 0 {
				t.Fatalf("missing tool modules: %s/%s", ch.ID, item.ID)
			}
			for _, answer := range append([]string{item.Answer.Primary}, item.Answer.Accepted...) {
				if !item.MatchAnswer(answer) {
					t.Fatalf("declared answer rejected: %s/%s %q", ch.ID, item.ID, answer)
				}
			}
			if item.MatchAnswer("") || item.MatchAnswer("definitely-wrong-answer") {
				t.Fatalf("false positive: %s/%s", ch.ID, item.ID)
			}
			families[ch.Family]++
			total++
		}
	}
	if total != 253 || !reflect.DeepEqual(families, map[Family]int{FamilyUnixShell: 184, FamilyKeyboardShortcuts: 31, FamilyModernCLIMacOS: 38}) {
		t.Fatalf("unexpected corpus: %d / %v", total, families)
	}
}

func TestV2MetadataIsExplicitInEverySource(t *testing.T) {
	paths, err := fs.Glob(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		raw, err := fs.ReadFile(data.Chapters(), path)
		if err != nil {
			t.Fatal(err)
		}
		var fields struct {
			Family        *string          `yaml:"family"`
			LearningOrder *string          `yaml:"learning_order"`
			Items         []map[string]any `yaml:"items"`
		}
		if err := yaml.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		if fields.Family == nil || fields.LearningOrder == nil {
			t.Fatalf("implicit chapter metadata: %s", path)
		}
		for _, item := range fields.Items {
			for _, name := range []string{"difficulty", "tools"} {
				if _, ok := item[name]; !ok {
					t.Fatalf("%s/%v: implicit %s", path, item["id"], name)
				}
			}
		}
	}
}

func TestV2HistoricalKeyPolicy(t *testing.T) {
	current, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	historical, err := LoadFS(os.DirFS("../../testdata/legacy-chapters"), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	old := map[string]bool{}
	for _, ch := range historical {
		for _, item := range ch.Items {
			old[ch.ID+"/"+item.ID] = true
		}
	}
	var preserved []string
	fresh := 0
	for _, ch := range current {
		for _, item := range ch.Items {
			key := ch.ID + "/" + item.ID
			if old[key] {
				preserved = append(preserved, key)
			} else {
				fresh++
			}
		}
	}
	want := []string{
		"01-raccourcis-terminal/ctrl-a", "01-raccourcis-terminal/ctrl-e", "01-raccourcis-terminal/ctrl-k", "01-raccourcis-terminal/ctrl-l", "01-raccourcis-terminal/ctrl-r", "01-raccourcis-terminal/ctrl-w", "01-raccourcis-terminal/ctrl-y",
		"02-navigation-shell/fd-find-dir", "02-navigation-shell/fd-fzf-dir", "05-outils-quotidiens/lazygit",
		"06-recherche-fichiers-contenu/fd-fzf-md", "06-recherche-fichiers-contenu/rg-content", "06-recherche-fichiers-contenu/rg-fzf-todo",
		"07-lecture-preview/bat-readme", "08-micro/micro-open-file", "08-micro/micro-quit", "08-micro/micro-save", "08-micro/micro-search",
		"09-markdown-glow/glow-readme", "10-workflows-dotfiles/dotfiles-status-shell",
	}
	sort.Strings(preserved)
	sort.Strings(want)
	if !reflect.DeepEqual(preserved, want) || fresh != 233 || len(old)-len(preserved) != 48 {
		t.Fatalf("key policy changed: kept=%v fresh=%d retired=%d", preserved, fresh, len(old)-len(preserved))
	}
}
