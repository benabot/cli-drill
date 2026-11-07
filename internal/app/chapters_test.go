package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/benabot/cli-drill/internal/progress"
	"github.com/benabot/cli-drill/internal/tui"
	"github.com/benabot/cli-drill/internal/xdg"
)

func TestChapterCommandsIgnoreUserYAML(t *testing.T) {
	for _, userYAML := range []struct {
		name    string
		content string
	}{
		{name: "valid override", content: embeddedTestChapter("User override", "echo user")},
		{name: "invalid YAML", content: "items: ["},
	} {
		for _, action := range []string{"chapters", "train", "train by id", "TUI"} {
			t.Run(userYAML.name+"/"+action, func(t *testing.T) {
				isolateChapterHome(t)
				userFiles := map[string]string{
					"override.yaml": userYAML.content,
					"extra.yaml": strings.ReplaceAll(
						embeddedTestChapter("User only", "echo extra"), "id: embedded", "id: user-only"),
				}
				for name, content := range userFiles {
					writeAppTestFile(t, xdg.ChaptersDir(), name, content)
				}

				var output bytes.Buffer
				cmd := NewRootCommand(fstest.MapFS{
					"embedded.yaml": {Data: []byte(embeddedTestChapter("Embedded lesson", "pwd"))},
				})
				cmd.SetOut(&output)
				cmd.SetErr(&output)
				cmd.SetIn(strings.NewReader("pwd\n"))
				// Training writes only this disposable session file, never the legacy JSON.
				sessionPath := filepath.Join(t.TempDir(), "session.json")
				args := []string{"--progress", sessionPath}
				switch action {
				case "chapters":
					args = append(args, "chapters")
				case "train":
					args = append(args, "train")
				case "train by id":
					args = append(args, "train", "embedded")
				case "TUI":
					// Keep chapter loading real; replace only the interactive terminal loop.
					restore := stubTUIRunner(func(opts tui.Options) (tui.Result, error) {
						if len(opts.Chapters) != 1 || opts.Chapters[0].ID != "embedded" ||
							opts.Chapters[0].Title != "Embedded lesson" ||
							len(opts.Chapters[0].Items) != 1 ||
							opts.Chapters[0].Items[0].Answer.Primary != "pwd" {
							t.Fatalf("TUI received user chapters: %#v", opts.Chapters)
						}
						return tui.Result{}, nil
					})
					t.Cleanup(restore)
				}
				cmd.SetArgs(args)
				if err := cmd.Execute(); err != nil {
					t.Fatalf("Execute: %v\n%s", err, output.String())
				}
				if action == "chapters" && output.String() != "embedded\tEmbedded lesson\t1 exercises\n" {
					t.Fatalf("unexpected chapter list: %s", output.String())
				}
				if strings.HasPrefix(action, "train") {
					if !strings.Contains(output.String(), "Chapitre: Embedded lesson") ||
						!strings.Contains(output.String(), "Correct.") {
						t.Fatalf("training did not use the embedded answer: %s", output.String())
					}
					state, err := progress.Load(sessionPath)
					if err != nil {
						t.Fatal(err)
					}
					if len(state.CompletedExercises) != 1 || state.CompletedExercises[0] != "embedded/pwd" {
						t.Fatalf("exercise identity changed: %#v", state.CompletedExercises)
					}
				}
				for name, original := range userFiles {
					got, err := os.ReadFile(filepath.Join(xdg.ChaptersDir(), name))
					if err != nil || string(got) != original {
						t.Errorf("user YAML %s changed: %q, %v", name, got, err)
					}
				}
				if _, err := os.Stat(xdg.ConfigFile()); !os.IsNotExist(err) {
					t.Fatalf("a configuration was created: %v", err)
				}
			})
		}
	}
}

func TestEmbeddedChapterLoadingPreservesLegacyData(t *testing.T) {
	isolateChapterHome(t)
	files := map[string]string{
		filepath.Join(xdg.ChaptersDir(), "user.yaml"): strings.ReplaceAll(
			embeddedTestChapter("User only", "echo user"), "id: embedded", "id: user-only"),
		xdg.ConfigFile():   "# Existing config must not be rewritten\nshell = \"zsh\"\n",
		xdg.ProgressFile(): `{"version":1,"completed_exercises":["embedded/pwd","user-only/old"],"chapter_scores":{"embedded":{"attempts":3,"correct":1},"user-only":{"attempts":4,"correct":2}},"last_session":"2026-06-06T12:00:00Z","streak":0}`,
	}
	for path, content := range files {
		writeAppTestFile(t, filepath.Dir(path), filepath.Base(path), content)
	}
	for _, action := range []string{"chapters", "stats", "train"} {
		var output bytes.Buffer
		cmd := NewRootCommand(fstest.MapFS{
			"embedded.yaml": {Data: []byte(embeddedTestChapter("Embedded lesson", "pwd"))},
		})
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetIn(strings.NewReader(""))
		args := []string{action}
		if action == "train" {
			args = append(args, "user-only")
		}
		cmd.SetArgs(args)
		err := cmd.Execute()
		if action == "train" {
			if err == nil || !strings.Contains(err.Error(), "chapter not found: user-only") {
				t.Fatalf("user-only chapter was available: %v\n%s", err, output.String())
			}
		} else if err != nil {
			t.Fatalf("%s: %v\n%s", action, err, output.String())
		}
		if action == "stats" && (!strings.Contains(output.String(), "Completed exercises: 2") ||
			!strings.Contains(output.String(), "Chapters with attempts: 2")) {
			t.Fatalf("legacy progress was not loaded: %s", output.String())
		}
	}
	for path, original := range files {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != original {
			t.Errorf("legacy data changed at %s: %q, %v", path, got, err)
		}
	}
}

func isolateChapterHome(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	// UserConfigDir uses HOME on macOS and XDG_CONFIG_HOME on Unix; isolate both.
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "config"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "data"))
	t.Setenv("AppData", filepath.Join(home, "config"))
}

func embeddedTestChapter(title, answer string) string {
	return "id: embedded\ntitle: " + title + "\nitems:\n" +
		"  - id: pwd\n    type: workflow\n    exercise_type: simple-shell-sim\n" +
		"    prompt: Print the current directory\n    answer:\n      primary: " + answer + "\n"
}
