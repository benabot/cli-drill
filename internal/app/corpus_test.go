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
	dir := useTemporaryTrainingDirectory(t)
	path := filepath.Join(dir, "progress.json")
	output, err := runSessionCLI(path, "pwd\nls\nls -a\ncd docs\ncd ..\ncd /atelier\nmkdir essais\ncat 'notes de cours.txt'\ncd ~\nman ls\n", "train")
	if err != nil || !strings.Contains(output, "Premiers pas Unix") || strings.Count(output, "Correct.") != 10 {
		t.Fatalf("default training failed: %v\n%s", err, output)
	}
	state, err := progress.Load(path)
	if err != nil || state.ChapterScores["00-premiers-pas"] != (progress.ChapterScore{Attempts: 10, Correct: 10}) || len(state.CompletedExercises) != 10 {
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

func TestWave3MixedTrainingRecordsKeysAndText(t *testing.T) {
	isolateChapterHome(t)
	allowSimulatedCapture(t)
	path := filepath.Join(t.TempDir(), "progress.json")
	const chapterID = "01-raccourcis-terminal"
	writeAppTestFileContent(t, path, `{"version":1,"completed_exercises":["retired/item","01-raccourcis-terminal/ctrl-a"],"chapter_scores":{"retired":{"attempts":7,"correct":3},"01-raccourcis-terminal":{"attempts":1,"correct":1}},"streak":4}`)
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ch, ok := findChapter(chapters, chapterID)
	if !ok {
		t.Fatal("missing line editing chapter")
	}
	items := ch.Items
	ch.Items = nil
	for _, id := range []string{"ctrl-d-delete-under-cursor", "readline-backward-word-edit", "ctrl-underscore-undo-edit"} {
		index := slices.IndexFunc(items, func(item chapter.Item) bool { return item.ID == id })
		if index < 0 {
			t.Fatalf("missing exercise %s", id)
		}
		ch.Items = append(ch.Items, items[index])
	}
	// This fixture feeds bytes to the mixed loop. The PTY smoke test exercises raw mode.
	var out bytes.Buffer
	input := strings.NewReader("\x04\ngit switch |sauvegarde\n\x1f\n")
	err = runTraining(Options{In: input, Out: &out, ProgressPath: path}, ch)
	if err != nil || strings.Count(out.String(), "Correct.") != 3 {
		t.Fatalf("mixed session: %v\n%s", err, out.String())
	}
	state := assertRecordingScore(t, path, chapterID, 4, 4, chapterID+"/ctrl-a",
		chapterID+"/ctrl-d-delete-under-cursor", chapterID+"/readline-backward-word-edit", chapterID+"/ctrl-underscore-undo-edit")
	if state.Version != 1 {
		t.Fatalf("unexpected progress version: %d", state.Version)
	}
}

func TestUnixStreamsRecordWithoutRunningRedirections(t *testing.T) {
	isolateChapterHome(t)
	dir := useTemporaryTrainingDirectory(t)
	path := filepath.Join(dir, "progress.json")
	writeAppTestFileContent(t, path, recordingHistory)
	prefs := writeAppTestFileContent(t, xdg.PreferencesFile(), `{"version":1,"enabled_tools":[],"default_difficulty":"beginner"}`)
	inventory := writeAppTestFileContent(t, "inventaire.txt", "inventaire précédent\n")
	journal := writeAppTestFileContent(t, "journal.txt", "entrée précédente\n")
	writeAppTestFileContent(t, "note.txt", "nouvelle entrée\n")
	writeAppTestFileContent(t, filepath.Join("docs", "fichier.txt"), "documentation\n")

	output, err := runSessionCLI(path, "ls docs > inventaire.txt\ncat note.txt > journal.txt\n", "train", "11-flux-shell", "--mode", "learn")
	if err != nil || !strings.Contains(output, "Bilan de la séance\nBonnes réponses : 1/2") {
		t.Fatalf("stream training: %v\n%s", err, output)
	}
	state := assertRecordingScore(t, path, "11-flux-shell", 2, 1, "11-flux-shell/redirect-output")
	if state.Version != 1 {
		t.Fatalf("unexpected progress version: %d", state.Version)
	}
	assertPreferencesBytes(t, xdg.PreferencesFile(), prefs)
	assertPreferencesBytes(t, "inventaire.txt", inventory)
	assertPreferencesBytes(t, "journal.txt", journal)
}

func TestTextTrainingPlainOutputSeparatesQuestionAndFeedback(t *testing.T) {
	isolateChapterHome(t)
	t.Setenv("CLICOLOR_FORCE", "1")
	path := filepath.Join(t.TempDir(), "progress.json")
	args := []string{"train", "11-flux-shell", "--difficulty", "advanced"}
	const solution = "ls docs absent > bilan.log 2>&1"
	question, err := runSessionCLI(path, "", args...)
	if err != nil || !strings.Contains(question, "Exercice 1/10\n\nListe le dossier docs") ||
		!strings.HasSuffix(question, "\n\nRéponse\n> ") || strings.Contains(question, solution) ||
		strings.Contains(question, "Explication") || strings.Contains(question, "\x1b") {
		t.Fatalf("unexpected unanswered exercise: %v\n%s", err, question)
	}
	for _, tc := range []struct{ answer, feedback, score string }{
		{solution, "Correct.\n\nExplication", "Bonnes réponses : 1/1"},
		{"ls docs absent 2>&1 > bilan.log", "Réponse incorrecte.\nRéponse attendue : " + solution + "\n\nExplication", "Bonnes réponses : 0/1"},
	} {
		// EOF after the answer stops the run before the next question.
		output, err := runSessionCLI(path, tc.answer, args...)
		if err != nil || !strings.HasPrefix(output, question+"\n"+tc.feedback+"\n") ||
			!strings.Contains(output, "Le shell redirige d’abord") ||
			!strings.HasSuffix(output, "Bilan de la séance\n"+tc.score+"\n") ||
			strings.Contains(output, "\x1b") {
			t.Fatalf("unexpected feedback: %v\n%s", err, output)
		}
	}
}

func useTemporaryTrainingDirectory(t *testing.T) string {
	t.Helper()
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
	return dir
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
