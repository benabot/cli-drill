package chapter

import (
	"slices"
	"testing"

	"github.com/benabot/cli-drill/data"
)

func TestModernExpansionEligibilityAndAnswers(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		chapter, id   string
		level         Difficulty
		tools         []string
		answer, wrong string
	}{
		{"02-navigation-shell", "eza-tree-depth", DifficultyIntermediate, []string{"eza"}, "eza --tree --level=2", "eza --tree"},
		{"02-navigation-shell", "eza-hidden-long", DifficultyIntermediate, []string{"eza"}, "eza -la", "eza -l"},
		{"02-navigation-shell", "zoxide-list-matches", DifficultyIntermediate, []string{"zoxide"}, "zoxide query -l docs", "zoxide query docs"},
		{"05-outils-quotidiens", "atuin-failed-search", DifficultyIntermediate, []string{"atuin"}, "atuin search --exit 1 git", "atuin search git"},
		{"06-recherche-fichiers-contenu", "rg-numbered", DifficultyIntermediate, []string{"rg"}, "rg -n TODO docs", "rg TODO docs"},
		{"06-recherche-fichiers-contenu", "rg-matching-files", DifficultyIntermediate, []string{"rg"}, "rg -l TODO docs", "rg TODO docs"},
		{"06-recherche-fichiers-contenu", "rg-fixed-dot", DifficultyIntermediate, []string{"rg"}, "rg -F 'v1.2' notes.txt", "rg v1.2 notes.txt"},
		{"06-recherche-fichiers-contenu", "rg-go-glob", DifficultyIntermediate, []string{"rg"}, "rg -g '*.go' TODO", "rg TODO"},
		{"06-recherche-fichiers-contenu", "fd-go-files", DifficultyIntermediate, []string{"fd"}, "fd -t f -e go", "fd -e go"},
		{"06-recherche-fichiers-contenu", "fd-hidden-notes", DifficultyIntermediate, []string{"fd"}, "fd -H -t f notes", "fd -t f notes"},
		{"06-recherche-fichiers-contenu", "rg-files-select", DifficultyAdvanced, []string{"rg", "fzf", "bat"}, `bat "$(rg -l TODO docs | fzf)"`, `bat $(rg -l TODO docs | fzf)`},
		{"07-lecture-preview", "bat-number-lines", DifficultyIntermediate, []string{"bat"}, "bat -n README.md", "bat README.md"},
		{"07-lecture-preview", "bat-line-range", DifficultyIntermediate, []string{"bat"}, "bat -r 10:20 notes.txt", "bat notes.txt"},
		{"07-lecture-preview", "bat-json-stdin", DifficultyAdvanced, []string{"bat"}, "bat --color=always --paging=never -l json", "bat"},
		{"08-micro", "micro-undo", DifficultyBeginner, []string{"micro"}, "Ctrl+Z", "Ctrl+Y"},
		{"08-micro", "micro-line-numbers", DifficultyIntermediate, []string{"micro"}, "Ctrl+R", "Ctrl+G"},
		{"09-markdown-glow", "glow-pager", DifficultyIntermediate, []string{"glow"}, "glow -p README.md", "glow README.md"},
		{"09-markdown-glow", "glow-width", DifficultyIntermediate, []string{"glow"}, "glow -w 60 README.md", "glow README.md"},
	}
	for _, tc := range cases {
		t.Run(tc.chapter+"/"+tc.id, func(t *testing.T) {
			idx := slices.IndexFunc(chapters, func(ch Chapter) bool { return ch.ID == tc.chapter })
			if idx < 0 {
				t.Fatal("chapter missing")
			}
			ch := chapters[idx]
			i := slices.IndexFunc(ch.Items, func(item Item) bool { return item.ID == tc.id })
			if i < 0 {
				t.Fatal("new exercise missing")
			}
			item := ch.Items[i]
			if item.Difficulty != tc.level || !slices.Equal(item.Tools, tc.tools) || !item.MatchAnswer(tc.answer) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("invalid metadata or correction: %#v", item)
			}
			filter := SessionFilter{Difficulty: tc.level, AvailableTools: AvailableTools(chapters)}
			if slices.ContainsFunc(filter.EligibleItems(ch, AllExercises), func(got Item) bool { return got.ID == tc.id }) {
				t.Fatal("exercise available without required modules")
			}
			filter.EnabledTools = tc.tools
			eligible := filter.EligibleItems(ch, AllExercises)
			if !slices.ContainsFunc(eligible, func(got Item) bool { return got.ID == tc.id }) {
				t.Fatal("exercise missing with required modules enabled")
			}
			if len(tc.tools) > 1 {
				for missing := range tc.tools {
					filter.EnabledTools = slices.Delete(slices.Clone(tc.tools), missing, missing+1)
					if slices.ContainsFunc(filter.EligibleItems(ch, AllExercises), func(got Item) bool { return got.ID == tc.id }) {
						t.Fatalf("multi-tool exercise available without %s", tc.tools[missing])
					}
				}
			}
			filter.EnabledTools = tc.tools
			shuffled := filter.Prepare(ch, AllExercises, "", func(n int, swap func(int, int)) {
				for i := 0; i < n/2; i++ {
					swap(i, n-1-i)
				}
			})
			slices.Reverse(eligible)
			if len(shuffled.Items) != len(eligible) {
				t.Fatal("shuffle lost exercises")
			}
			for i := range eligible {
				if eligible[i].ID != shuffled.Items[i].ID {
					t.Fatal("shuffle changed the eligible set")
				}
			}
		})
	}
}
