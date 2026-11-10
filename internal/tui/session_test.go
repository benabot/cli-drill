package tui

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func TestChapterDetailSessionChoicesAndStableRetry(t *testing.T) {
	ch := recordingChapter()
	ch.Items[1].Difficulty = chapter.DifficultyIntermediate
	ch.Items[1].Tools = []string{"fd"}
	ch.Items[2].Difficulty = chapter.DifficultyIntermediate
	ch.Items[2].Tools = []string{"fd", "fzf"}
	ch.Items[3].Difficulty = chapter.DifficultyAdvanced
	path := recordingPath(t)
	level := chapter.DifficultyBeginner
	tools := []string{"fd", "fzf", "retired-tool"}
	opts := Options{
		Chapters: []chapter.Chapter{ch}, ProgressPath: path,
		LoadPreferences: func() ([]string, chapter.Difficulty, error) { return tools, level, nil },
		SavePreferences: func([]string, chapter.Difficulty) error {
			t.Fatal("session choices wrote preferences")
			return nil
		},
	}
	m := newModel(opts)
	m.home.Select(homeActionIndex(homeActionBrowse))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Session difficulty: beginner")
	assertContains(t, m.View(), "Mode: Learn (editorial order)")
	m = press(t, m, tea.KeyRunes, 'd')
	m = press(t, m, tea.KeyRunes, 'm')
	m = press(t, m, tea.KeyRunes, 'h')
	assertContains(t, m.View(), "Session difficulty: intermediate")
	assertContains(t, m.View(), "Mode: Review (shuffled)")
	assertContains(t, m.View(), "intermediate (2)")
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	order := tuiSessionIDs(m.session.Items)
	sorted := slices.Clone(order)
	slices.Sort(sorted)
	if !slices.Equal(sorted, []string{"multiple-choice", "scenario"}) || m.showHelp {
		t.Fatalf("wrong prepared session: %v", order)
	}
	// Changing preferences during a session leaves the active slice intact.
	level, tools = chapter.DifficultyAdvanced, nil

	target := path + ".target"
	if err := os.Rename(path, target); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	if !m.isMultipleChoice() {
		m = assertSessionKeysReachTextInput(t, m)
		m.input.SetValue("pwd")
	}
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Progress could not be saved")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(model)
	if m.itemIndex != 0 || m.summary.Total != 0 || !slices.Equal(order, tuiSessionIDs(m.session.Items)) {
		t.Fatal("save failure or resize changed the active session")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(target, path); err != nil {
		t.Fatal(err)
	}
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Correct.")
	if !slices.Equal(order, tuiSessionIDs(m.session.Items)) {
		t.Fatal("feedback changed order")
	}
	m = press(t, m, tea.KeyEnter)
	if m.currentItem().ID != order[1] {
		t.Fatalf("advanced to %s, want %s", m.currentItem().ID, order[1])
	}
	if !m.isMultipleChoice() {
		m = assertSessionKeysReachTextInput(t, m)
	}
	assertSavedRecording(t, path, 1, 1, ch.ID+"/"+order[0])
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenChapterDetail)
	assertContains(t, m.View(), "Session difficulty: advanced")
	assertContains(t, m.View(), "Mode: Learn (editorial order)")
	m = press(t, m, tea.KeyEnter)
	if !slices.Equal(tuiSessionIDs(m.session.Items), []string{"simple-shell-sim"}) {
		t.Fatalf("next session ignored updated preferences: %v", tuiSessionIDs(m.session.Items))
	}
}

func assertSessionKeysReachTextInput(t *testing.T, m model) model {
	t.Helper()
	for _, key := range "dmh" {
		m = press(t, m, tea.KeyRunes, key)
	}
	if m.input.Value() != "dmh" {
		t.Fatalf("session shortcuts intercepted text: %q", m.input.Value())
	}
	return m
}

func TestTUIEmptySelectionAndHistoricalAcquired(t *testing.T) {
	ch := recordingChapter()
	ch.Items = ch.Items[:2]
	ch.Items[1].Difficulty = chapter.DifficultyIntermediate
	ch.Items[1].Tools = []string{"fd", "fzf"}
	level := chapter.DifficultyIntermediate
	tools := []string{"fd", "retired"}
	path := filepath.Join(t.TempDir(), "progress.json")
	history := `{"version":1,"completed_exercises":["text/free-answer","text/scenario","text/removed"],"chapter_scores":{"text":{"attempts":9,"correct":5}},"last_session":null,"streak":4}`
	if err := os.WriteFile(path, []byte(history), 0o600); err != nil {
		t.Fatal(err)
	}
	m := newModel(Options{Chapters: []chapter.Chapter{ch}, ProgressPath: path,
		LoadPreferences: func() ([]string, chapter.Difficulty, error) { return tools, level, nil },
	})
	m.home.Select(homeActionIndex(homeActionStartTraining))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapters)
	assertContains(t, m.View(), "No eligible text exercises")
	assertContains(t, m.View(), "Exercices éligibles : 0/2")
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "with the enabled CLI modules")
	assertContains(t, m.View(), "Historical chapter score: 5/9 correct")
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapterDetail)
	m = press(t, m, tea.KeyRunes, 'd')
	assertContains(t, m.View(), "No exercises at advanced difficulty")
	m = press(t, m, tea.KeyRunes, 'd')
	assertContains(t, m.View(), "Active acquired: 1/1 eligible")
	m = press(t, m, tea.KeyRunes, 'g')
	m.home.Select(homeActionIndex(homeActionStats))
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Active acquired (intermediate, enabled CLI modules): 0/0")
	assertContains(t, m.View(), "Historical acquired: 3")
	assertContains(t, m.View(), "Chapters with historical scores: 1")
	m = press(t, m, tea.KeyEsc)
	tools = append(tools, "fzf")
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Active acquired (intermediate, enabled CLI modules): 1/1")
	got, err := os.ReadFile(path)
	if err != nil || string(got) != history {
		t.Fatalf("consultation or empty launch modified progress: %s (%v)", got, err)
	}
}

func TestKeyboardShortcutsOnlyDelegatesEligibleKeys(t *testing.T) {
	legacy := keySequenceChapter()
	mixed := recordingChapter()
	mixed.Items[0].Difficulty = chapter.DifficultyIntermediate
	key := legacy.Items[0]
	key.Difficulty = chapter.DifficultyIntermediate
	key.Tools = []string{"fd", "fzf"}
	mixed.Items = append(mixed.Items, key)
	level := chapter.DifficultyIntermediate
	tools := []string{"fd"}
	path := recordingPath(t)
	opts := Options{Chapters: []chapter.Chapter{legacy, mixed}, ProgressPath: path,
		LoadPreferences: func() ([]string, chapter.Difficulty, error) { return tools, level, nil },
	}
	m := newModel(opts)
	m.home.Select(homeActionIndex(homeActionShortcuts))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapters)
	assertContains(t, m.View(), "No eligible keyboard shortcuts")
	m = press(t, m, tea.KeyEsc)

	tools = append(tools, "fzf")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	capture := updated.(model).result.CaptureChapter
	assertQuitCommand(t, cmd)
	if capture.ID != mixed.ID || !slices.Equal(tuiSessionIDs(capture.Items), []string{key.ID}) {
		t.Fatalf("keyboard entry delegated wrong items: %#v", capture)
	}
	state, err := progress.Load(path)
	if err != nil || len(state.ChapterScores) != 1 || len(state.CompletedExercises) != 1 {
		t.Fatalf("preparing capture wrote progress: %#v (%v)", state, err)
	}

	opts.CaptureUnavailable = "Keyboard capture requires a terminal."
	m = newModel(opts)
	m.home.Select(homeActionIndex(homeActionShortcuts))
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd != nil || m.result.CaptureChapter.ID != "" {
		t.Fatal("started capture without terminal support")
	}
	assertContains(t, m.View(), opts.CaptureUnavailable)
}

func tuiSessionIDs(items []chapter.Item) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
