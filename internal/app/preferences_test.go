package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/xdg"
)

func TestPreferencesLoadAndExplicitSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cli-drill", "preferences.json")
	initial, err := loadPreferences(path)
	if err != nil || initial.Version != 1 || initial.EnabledTools == nil || len(initial.EnabledTools) != 0 ||
		initial.DefaultDifficulty != chapter.DifficultyBeginner {
		t.Fatalf("defaults = %+v (%v)", initial, err)
	}
	if _, err := os.Stat(filepath.Dir(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("load created storage: %v", err)
	}

	updated := initial
	updated.EnabledTools = []string{"fd", "retired", "fzf"}
	updated.DefaultDifficulty = chapter.DifficultyAdvanced
	if err := savePreferences(path, initial, updated); err != nil {
		t.Fatal(err)
	}
	reloaded, err := loadPreferences(path)
	if err != nil || !slices.Equal(reloaded.EnabledTools, updated.EnabledTools) || reloaded.DefaultDifficulty != updated.DefaultDifficulty {
		t.Fatalf("reloaded = %+v (%v)", reloaded, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("preferences are not private: %v", info.Mode())
	}
	if err := savePreferences(path, initial, initial); err == nil {
		t.Fatal("stale editor overwrote newer preferences")
	}
}

func TestPreferencesOptionalFieldsAndInvalidFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preferences.json")
	write := func(contents string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write(`{"version":1}`)
	state, err := loadPreferences(path)
	if err != nil || state.EnabledTools == nil || len(state.EnabledTools) != 0 || state.DefaultDifficulty != chapter.DifficultyBeginner {
		t.Fatalf("omitted fields = %+v (%v)", state, err)
	}
	for _, contents := range []string{
		`{`, `null`, `[]`, `{}`, `{"version":null}`, `{"version":2}`, `{"version":"1"}`,
		`{"version":1,"enabled_tools":null}`, `{"version":1,"enabled_tools":"fd"}`,
		`{"version":1,"enabled_tools":[1]}`, `{"version":1,"enabled_tools":[null]}`,
		`{"version":1,"enabled_tools":[" "]}`, `{"version":1,"default_difficulty":null}`,
		`{"version":1,"default_difficulty":2}`, `{"version":1,"default_difficulty":"expert"}`,
	} {
		t.Run(contents, func(t *testing.T) {
			write(contents)
			if _, err := loadPreferences(path); err == nil {
				t.Fatal("invalid preferences loaded")
			}
			if err := savePreferences(path, state, state); err == nil {
				t.Fatal("invalid file overwritten")
			}
			assertPreferencesBytes(t, path, []byte(contents))
		})
	}
}

func TestPreferencesCLIAndProgressIsolation(t *testing.T) {
	isolateChapterHome(t)
	path := xdg.PreferencesFile()
	override := filepath.Join(t.TempDir(), "override.json")
	run := func(args ...string) (string, error) {
		return runPreferencesCLI(t, append([]string{"--progress", override}, args...)...)
	}

	out, err := run("prefs")
	if err != nil || !strings.Contains(out, "Default difficulty: beginner") || !strings.Contains(out, "fd: disabled") {
		t.Fatalf("defaults: %q (%v)", out, err)
	}
	if out, err := run("prefs", "--disable", "fd"); err != nil || !strings.Contains(out, "Preferences unchanged.") {
		t.Fatalf("no-op mutation: %q (%v)", out, err)
	}
	for _, filename := range []string{path, override, xdg.ProgressFile()} {
		if _, err := os.Stat(filename); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("consultation created %s (%v)", filename, err)
		}
	}

	history := []byte(`{"version":1,"completed_exercises":["old/item"],"chapter_scores":{"old":{"attempts":2,"correct":1}},"last_session":null,"streak":0}`)
	for _, filename := range []string{override, xdg.ProgressFile()} {
		writeAppTestFile(t, filepath.Dir(filename), filepath.Base(filename), string(history))
	}
	writeAppTestFile(t, filepath.Dir(path), filepath.Base(path), `{"version":1,"enabled_tools":["retired-cli"]}`)
	legacy := map[string]string{
		filepath.Join(legacyConfigRoot(t), "config.toml"):          "[invalid TOML",
		filepath.Join(legacyConfigRoot(t), "chapters", "old.yaml"): "items: [",
	}
	for filename, contents := range legacy {
		writeAppTestFile(t, filepath.Dir(filename), filepath.Base(filename), contents)
	}

	for _, args := range [][]string{
		{"prefs", "--enable", "fd"},
		{"prefs", "--enable", "fzf"},
		{"prefs", "--disable", "fd"},
		{"prefs", "--difficulty", "intermediate"},
	} {
		if out, err := run(args...); err != nil || !strings.Contains(out, "Preferences saved.") {
			t.Fatalf("%v: %q (%v)", args, out, err)
		}
	}
	out, err = run("prefs")
	for _, want := range []string{"fd: disabled", "fzf: enabled", "Default difficulty: intermediate", "retired-cli"} {
		if err != nil || !strings.Contains(out, want) {
			t.Fatalf("missing %q: %q (%v)", want, out, err)
		}
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"prefs", "--enable", "retired-cli"},
		{"prefs", "--disable", "retired-cli"},
		{"prefs", "--difficulty", "expert"},
		{"prefs", "--enable", "fd", "--disable", "fzf"},
	} {
		if out, err := run(args...); err == nil || strings.Contains(out, "Preferences saved.") {
			t.Fatalf("invalid invocation %v: %q (%v)", args, out, err)
		}
		assertPreferencesBytes(t, path, saved)
	}
	if out, err := run("stats"); err != nil || !strings.Contains(out, "Completed exercises: 1") {
		t.Fatalf("historical progress unreadable: %q (%v)", out, err)
	}
	assertPreferencesBytes(t, override, history)
	assertPreferencesBytes(t, xdg.ProgressFile(), history)
	if out, err := run("reset", "--yes"); err != nil {
		t.Fatalf("reset: %q (%v)", out, err)
	}
	if _, err := os.Stat(override); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("reset did not remove override: %v", err)
	}
	assertPreferencesBytes(t, xdg.ProgressFile(), history)
	assertPreferencesBytes(t, path, saved)
	for filename, original := range legacy {
		assertPreferencesBytes(t, filename, []byte(original))
	}
}

func TestProgressOverrideRejectsPreferencesAndAliases(t *testing.T) {
	isolateChapterHome(t)
	path := xdg.PreferencesFile()
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(parent, alias); err != nil {
		t.Fatal(err)
	}
	viaAlias := filepath.Join(alias, filepath.Base(path))
	// Reject the reserved target even before it exists, so train cannot create
	// progress at the preferences location.
	if _, err := runPreferencesCLI(t, "--progress", viaAlias, "train"); err == nil {
		t.Fatal("train accepted the reserved file through an aliased directory")
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("reserved file unexpectedly exists: %v", err)
	}
	original := []byte(`{"version":1,"enabled_tools":["fd"]}`)
	writeAppTestFile(t, parent, filepath.Base(path), string(original))
	symlink := filepath.Join(t.TempDir(), "alternate.json")
	if err := os.Symlink(path, symlink); err != nil {
		t.Fatal(err)
	}
	for _, override := range []string{path, viaAlias, symlink} {
		for _, action := range [][]string{{"reset", "--yes"}, {"train"}} {
			out, err := runPreferencesCLI(t, append([]string{"--progress", override}, action...)...)
			if err == nil || !strings.Contains(out, "preferences") {
				t.Fatalf("--progress points at preferences: %q (%v)", out, err)
			}
			assertPreferencesBytes(t, path, original)
		}
	}
}

func TestPreferencesCLIReportsUnusableStorage(t *testing.T) {
	isolateChapterHome(t)
	path := xdg.PreferencesFile()
	contents := `{"version":2}`
	writeAppTestFile(t, filepath.Dir(path), filepath.Base(path), contents)
	out, err := runPreferencesCLI(t, "prefs")
	if err != nil || !strings.Contains(out, "defaults") || !strings.Contains(out, "version") || !strings.Contains(out, "beginner") {
		t.Fatalf("invalid file diagnostic: %q (%v)", out, err)
	}
	out, err = runPreferencesCLI(t, "prefs", "--enable", "fd")
	if err == nil || strings.Contains(out, "Preferences saved.") {
		t.Fatalf("write against invalid file: %q (%v)", out, err)
	}
	assertPreferencesBytes(t, path, []byte(contents))

	// A broken parent symlink allows the missing-file read but makes creation fail.
	blocked := filepath.Join(t.TempDir(), "missing")
	if err := os.Symlink(blocked, filepath.Join(filepath.Dir(blocked), "cli-drill")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_DATA_HOME", filepath.Dir(blocked))
	out, err = runPreferencesCLI(t, "prefs", "--enable", "fd")
	if err == nil || !strings.Contains(err.Error(), "save preferences") || strings.Contains(out, "Preferences saved.") {
		t.Fatalf("failed write: %q (%v)", out, err)
	}
	if _, err := os.Stat(blocked); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("failed write created storage: %v", err)
	}
}

func runPreferencesCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetIn(strings.NewReader(""))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func assertPreferencesBytes(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("%s changed: %q (%v)", path, got, err)
	}
}
