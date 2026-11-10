package app

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/progress"
	"github.com/benabot/cli-drill/internal/xdg"
)

func TestDefaultTrainingUsesUnixFoundationsWithoutExecutingAnswers(t *testing.T) {
	isolateChapterHome(t)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatal(err)
		}
	})
	path := filepath.Join(dir, "progress.json")
	var output bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetIn(strings.NewReader("pwd\nls\nls -a\ncd docs\ncd ..\ncd /atelier\nmkdir essais\ncat 'notes de cours.txt'\n"))
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"--progress", path, "train"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("train: %v\n%s", err, output.String())
	}
	if !strings.Contains(output.String(), "Premiers pas Unix") || strings.Count(output.String(), "Correct.") != 8 {
		t.Fatalf("default training failed: %s", output.String())
	}
	state, err := progress.Load(path)
	if err != nil || state.ChapterScores["00-premiers-pas"] != (progress.ChapterScore{Attempts: 8, Correct: 8}) || len(state.CompletedExercises) != 8 {
		t.Fatalf("first training progress: %#v (%v)", state, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 || entries[0].Name() != "progress.json" {
		t.Fatalf("answers created files/directories: %v (%v)", entries, err)
	}
	physicalDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := os.Getwd(); err != nil || got != physicalDir {
		t.Fatalf("answer changed directory: %s (%v)", got, err)
	}
}

func TestCurrentTextCorpusPreservesRetiredProgress(t *testing.T) {
	isolateChapterHome(t)
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	prefs := defaultPreferences()
	prefs.EnabledTools = chapter.AvailableTools(chapters)
	dataJSON, err := json.Marshal(prefs)
	if err != nil {
		t.Fatal(err)
	}
	writeAppTestFileContent(t, xdg.PreferencesFile(), string(dataJSON))
	for _, ch := range chapters {
		if ch.HasKeyExercises() {
			continue
		}
		t.Run(ch.ID, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			const historical = `{"version":1,"completed_exercises":["03-alias-zsh/cgs","07-lecture-preview/marked2-preview","01-raccourcis-terminal/ctrl-u"],"chapter_scores":{"03-alias-zsh":{"attempts":9,"correct":4}},"last_session":null,"streak":4}`
			writeAppTestFileContent(t, path, historical)
			want := []string{"03-alias-zsh/cgs", "07-lecture-preview/marked2-preview", "01-raccourcis-terminal/ctrl-u"}
			for _, level := range []chapter.Difficulty{chapter.DifficultyBeginner, chapter.DifficultyIntermediate, chapter.DifficultyAdvanced} {
				var items []chapter.Item
				for _, item := range ch.Items {
					if item.Difficulty == level {
						items = append(items, item)
						want = append(want, ch.ID+"/"+item.ID)
					}
				}
				if len(items) == 0 {
					continue
				}
				var output bytes.Buffer
				lastRead := 0
				cmd := NewRootCommand(data.Chapters())
				cmd.SetIn(recordingReader(func(buffer []byte) (int, error) {
					rendered := output.String()[lastRead:]
					lastRead = output.Len()
					for _, item := range items {
						if strings.Contains(rendered, item.Prompt) {
							return copy(buffer, item.Answer.Primary+"\n"), nil
						}
					}
					t.Errorf("unexpected prompt at %s: %s", level, rendered)
					return 0, io.EOF
				}))
				cmd.SetOut(&output)
				cmd.SetErr(&output)
				cmd.SetArgs([]string{"--progress", path, "train", ch.ID, "--difficulty", string(level)})
				if err := cmd.Execute(); err != nil {
					t.Fatal(err)
				}
			}
			state, err := progress.Load(path)
			if err != nil {
				t.Fatal(err)
			}
			slices.Sort(want)
			got := slices.Clone(state.CompletedExercises)
			slices.Sort(got)
			if !slices.Equal(got, want) || state.ChapterScores[ch.ID] != (progress.ChapterScore{Attempts: len(ch.Items), Correct: len(ch.Items)}) || state.ChapterScores["03-alias-zsh"] != (progress.ChapterScore{Attempts: 9, Correct: 4}) || state.Version != 1 || state.Streak != 4 {
				t.Fatalf("current corpus lost history: %#v (want acquired %v)", state, want)
			}
		})
	}
}
