package chapter

import (
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/benabot/cli-drill/internal/exercise"
	"gopkg.in/yaml.v3"
)

func metadataYAML(chapterFields, itemFields string) string {
	return "id: legacy\ntitle: Legacy\n" + chapterFields + "items:\n" +
		"  - id: pwd\n    type: workflow\n    exercise_type: simple-shell-sim\n" +
		"    prompt: Print the current directory\n    answer:\n      primary: pwd\n" + itemFields
}

func TestLoadRejectsInvalidMetadata(t *testing.T) {
	for _, tc := range []struct {
		name, chapterFields, itemFields, field string
	}{
		{"difficulty unknown", "", "    difficulty: expert\n", "difficulty"},
		{"difficulty case", "", "    difficulty: Beginner\n", "difficulty"},
		{"difficulty empty", "", "    difficulty: ''\n", "difficulty"},
		{"difficulty null", "", "    difficulty: null\n", "difficulty"},
		{"difficulty sequence", "", "    difficulty: [beginner]\n", "difficulty"},
		{"family unknown", "family: other\n", "", "family"},
		{"family empty", "family: ''\n", "", "family"},
		{"family null", "family: null\n", "", "family"},
		{"family boolean", "family: true\n", "", "family"},
		{"order unknown", "learning_order: shuffled\n", "", "learning_order"},
		{"order empty", "learning_order: ''\n", "", "learning_order"},
		{"order null", "learning_order: null\n", "", "learning_order"},
		{"tools scalar", "", "    tools: rg\n", "tools"},
		{"tools null", "", "    tools: null\n", "tools"},
		{"tools map", "", "    tools: {rg: true}\n", "tools"},
		{"tool empty", "", "    tools: ['']\n", "tools"},
		{"tool blank", "", "    tools: ['  ']\n", "tools"},
		{"tool duplicate", "", "    tools: [fd, fd]\n", "tools"},
		{"tool number", "", "    tools: [42]\n", "tools"},
		{"tool boolean", "", "    tools: [true]\n", "tools"},
		{"tool null", "", "    tools: [null]\n", "tools"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadFS(fstest.MapFS{"lesson.yaml": {Data: []byte(metadataYAML(tc.chapterFields, tc.itemFields))}}, "*.yaml")
			if err == nil {
				t.Fatalf("accepted invalid %s", tc.field)
			}
			for _, want := range []string{"lesson.yaml", tc.field} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not identify %q", err, want)
				}
			}
		})
	}
}

func TestLoadMetadataDefaultsAndValues(t *testing.T) {
	legacy := metadataYAML("", "    source: {path: old.zsh}\n    tags: [historical]\n")
	ch, err := Load(strings.NewReader(legacy))
	if err != nil {
		t.Fatal(err)
	}
	if ch.Family != "" || ch.LearningOrder != LearningOrderGuided || ch.Items[0].Difficulty != DifficultyBeginner || len(ch.Items[0].Tools) != 0 {
		t.Fatalf("unexpected historical defaults: %#v / %#v", ch, ch.Items[0])
	}
	if ch.ID != "legacy" || ch.Items[0].ID != "pwd" {
		t.Fatal("historical identifiers changed")
	}
	for _, difficulty := range []Difficulty{DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced} {
		for _, family := range []Family{FamilyUnixShell, FamilyKeyboardShortcuts, FamilyModernCLIMacOS} {
			for _, order := range []LearningOrder{LearningOrderGuided, LearningOrderIndependent} {
				t.Run(string(difficulty)+"/"+string(family)+"/"+string(order), func(t *testing.T) {
					ch, err := Load(strings.NewReader(metadataYAML("family: "+string(family)+"\nlearning_order: "+string(order)+"\n", "    difficulty: "+string(difficulty)+"\n    tools: [fd, fzf]\n")))
					if err != nil {
						t.Fatal(err)
					}
					if ch.Family != family || ch.LearningOrder != order || ch.Items[0].Difficulty != difficulty || !reflect.DeepEqual(ch.Items[0].Tools, []string{"fd", "fzf"}) {
						t.Fatalf("metadata lost: %#v / %#v", ch, ch.Items[0])
					}
					if ch.ID != "legacy" || ch.Items[0].ID != "pwd" || !ch.Items[0].MatchAnswer("pwd") || ch.Items[0].MatchAnswer("PWD") {
						t.Fatal("identity or correction changed")
					}
				})
			}
		}
	}
	for _, fields := range []string{"", "    tools: []\n"} {
		ch, err := Load(strings.NewReader(metadataYAML("", fields)))
		if err != nil || len(ch.Items[0].Tools) != 0 {
			t.Fatalf("fundamentals should require no third-party module: %#v, %v", ch, err)
		}
	}
}

func TestLoadMetadataAliasesAndMerges(t *testing.T) {
	for _, fields := range []string{
		"    difficulty: &level advanced\n    tools: &modules [fd, fzf]\n",
		"    <<: {difficulty: advanced, tools: [fd, fzf]}\n",
		"    <<: {difficulty: beginner}\n    difficulty: advanced\n    tools: [fd, fzf]\n",
	} {
		ch, err := Load(strings.NewReader(metadataYAML("", fields)))
		if err != nil || ch.Items[0].Difficulty != DifficultyAdvanced || !reflect.DeepEqual(ch.Items[0].Tools, []string{"fd", "fzf"}) {
			t.Fatalf("YAML metadata lost: %#v (%v)", ch, err)
		}
	}
	input := metadataYAML("level: &level advanced\nmodules: &modules [fd, fzf]\n", "    difficulty: *level\n    tools: *modules\n")
	ch, err := Load(strings.NewReader(input))
	if err != nil || ch.Items[0].Difficulty != DifficultyAdvanced || len(ch.Items[0].Tools) != 2 {
		t.Fatalf("aliases were not resolved: %#v (%v)", ch, err)
	}
	for _, fields := range []string{"    <<: {difficulty: null}\n", "    <<: {tools: [null]}\n"} {
		if _, err := Load(strings.NewReader(metadataYAML("", fields))); err == nil {
			t.Fatal("merge bypassed metadata validation")
		}
	}
}

func TestLoadRejectsInvalidStructure(t *testing.T) {
	valid := metadataYAML("", "")
	for _, tc := range []struct{ name, input, want string }{
		{"exercise type", strings.ReplaceAll(valid, "simple-shell-sim", "unknown"), "exercise_type"},
		{"duplicate item", valid + strings.SplitN(valid, "items:\n", 2)[1], "duplicate item id"},
		{"duplicate metadata field", metadataYAML("family: unix-shell\nfamily: keyboard-shortcuts\n", ""), "family"},
		{"qcm missing choices", strings.ReplaceAll(valid, "simple-shell-sim", "multiple-choice"), "choices"},
		{"qcm no accepted choice", strings.ReplaceAll(metadataYAML("", "    choices: [wrong, unrelated]\n"), "simple-shell-sim", "multiple-choice"), "choices"},
		{"qcm duplicate choice", strings.ReplaceAll(metadataYAML("", "    choices: [pwd, pwd]\n"), "simple-shell-sim", "multiple-choice"), "choices"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(strings.NewReader(tc.input))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want error identifying %s", err, tc.want)
			}
		})
	}
	_, err := LoadFS(fstest.MapFS{"a.yaml": {Data: []byte(valid)}, "b.yaml": {Data: []byte(valid)}}, "*.yaml")
	if err == nil || !strings.Contains(err.Error(), "b.yaml") || !strings.Contains(err.Error(), "duplicate chapter id") {
		t.Fatalf("duplicate chapter not identified: %v", err)
	}
}

func TestMetadataPreservesHistoricalCorpus(t *testing.T) {
	source := os.DirFS("../../testdata/legacy-chapters")
	chapters, err := LoadFS(source, "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	paths, err := fs.Glob(source, "*.yaml")
	if err != nil || len(paths) != 10 || len(chapters) != 10 {
		t.Fatalf("embedded chapter count changed: %d (%v)", len(chapters), err)
	}
	total := 0
	for i, path := range paths {
		raw, err := fs.ReadFile(source, path)
		if err != nil {
			t.Fatal(err)
		}
		// Mirror the original model without custom decoding or P2 metadata.
		type legacyItem struct {
			ID           string          `yaml:"id"`
			Type         string          `yaml:"type"`
			ExerciseType exercise.Type   `yaml:"exercise_type"`
			Prompt       string          `yaml:"prompt"`
			Choices      []string        `yaml:"choices"`
			Answer       exercise.Answer `yaml:"answer"`
			Explanation  string          `yaml:"explanation"`
		}
		var historical struct {
			ID          string       `yaml:"id"`
			Title       string       `yaml:"title"`
			Description string       `yaml:"description"`
			Items       []legacyItem `yaml:"items"`
		}
		if err := yaml.Unmarshal(raw, &historical); err != nil {
			t.Fatal(err)
		}
		ch := chapters[i]
		if ch.ID != historical.ID || ch.Title != historical.Title || ch.Description != historical.Description || ch.Family != "" || ch.LearningOrder != LearningOrderGuided || len(ch.Items) != len(historical.Items) {
			t.Fatalf("historical chapter changed: %#v", ch)
		}
		for j, item := range ch.Items {
			got := legacyItem{ID: item.ID, Type: item.Type, ExerciseType: item.ExerciseType,
				Prompt: item.Prompt, Choices: item.Choices, Answer: item.Answer, Explanation: item.Explanation}
			if !reflect.DeepEqual(got, historical.Items[j]) || item.Difficulty != DifficultyBeginner || len(item.Tools) != 0 {
				t.Fatalf("historical item changed: %#v", item)
			}
			total++
		}
	}
	if total != 68 {
		t.Fatalf("got %d historical exercises, want 68", total)
	}
}
