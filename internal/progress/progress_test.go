package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLoadMissingFileDoesNotWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "progress.json")
	got, err := Load(path)
	if err != nil || !reflect.DeepEqual(got, Default()) {
		t.Fatalf("Load missing file = %#v, %v", got, err)
	}
	if _, err := os.Stat(filepath.Dir(path)); !os.IsNotExist(err) {
		t.Fatalf("Load created a directory: %v", err)
	}
}

func TestSaveAndLoadProgress(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := Default()
	progress.MarkCompleted("01-raccourcis-terminal", "ctrl-a", true)

	if err := Save(path, progress); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if len(loaded.CompletedExercises) != 1 {
		t.Fatalf("expected 1 completed exercise, got %#v", loaded.CompletedExercises)
	}
	if loaded.ChapterScores["01-raccourcis-terminal"].Correct != 1 {
		t.Fatalf("unexpected score: %#v", loaded.ChapterScores)
	}
}

func TestLoadHistoricalDefaults(t *testing.T) {
	for _, input := range []string{`{}`, `{"version":1}`, `{"version":1,"chapter_scores":null}`} {
		t.Run(input, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			writeFixture(t, path, []byte(input))
			got, err := Load(path)
			if err != nil || !reflect.DeepEqual(got, Default()) {
				t.Fatalf("Load sparse historical JSON = %#v, %v", got, err)
			}
			assertFileBytes(t, path, []byte(input))
		})
	}
}

func TestHistoricalProgressRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	when := time.Date(2026, 6, 6, 12, 30, 0, 0, time.UTC)
	want := Progress{
		Version:            1,
		CompletedExercises: []string{"retired-chapter/old-item", "01-raccourcis-terminal/ctrl-a"},
		ChapterScores: map[string]ChapterScore{
			"retired-chapter":        {Attempts: 7, Correct: 3},
			"01-raccourcis-terminal": {Attempts: 2, Correct: 1},
		},
		LastSession: &when,
		Streak:      4,
	}
	// Literal version-1 data, independent of the current serializer.
	original := []byte(`{"version":1,"completed_exercises":["retired-chapter/old-item","01-raccourcis-terminal/ctrl-a"],"chapter_scores":{"retired-chapter":{"attempts":7,"correct":3},"01-raccourcis-terminal":{"attempts":2,"correct":1}},"last_session":"2026-06-06T12:30:00Z","streak":4}`)
	writeFixture(t, path, original)
	loaded, err := Load(path)
	if err != nil || !reflect.DeepEqual(loaded, want) {
		t.Fatalf("historical Load = %#v, %v", loaded, err)
	}
	assertFileBytes(t, path, original)
	if err := Save(path, loaded); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip lost historical fields: %#v, %v", got, err)
	}
	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(written, &fields); err != nil || len(fields) != 5 {
		t.Fatalf("JSON schema changed: %s (%v)", written, err)
	}
	if !bytes.HasSuffix(written, []byte("\n")) {
		t.Fatal("missing historical trailing newline")
	}
	assertNoTemporaryFiles(t, filepath.Dir(path), "progress.json")
}

func TestInvalidProgressIsNeverOverwritten(t *testing.T) {
	for _, input := range []string{"", "{", "null", "[]", `{"version":2}`, `{"version":0}`, `{"version":1,"chapter_scores":"wrong type"}`, `{"version":1} {}`} {
		t.Run(input, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "progress.json")
			writeFixture(t, path, []byte(input))
			if _, err := Load(path); err == nil || !strings.Contains(err.Error(), path) {
				t.Errorf("Load must return a contextual error: %v", err)
			}
			if err := Save(path, Default()); err == nil {
				t.Error("Save must refuse an invalid existing file")
			}
			assertFileBytes(t, path, []byte(input))
			assertNoTemporaryFiles(t, dir, "progress.json")
		})
	}
}

func TestReadErrorDoesNotRemovePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected a read error for a directory")
	}
	if err := Save(path, Default()); err == nil {
		t.Fatal("expected a save error for a directory")
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		t.Fatalf("directory was removed: %v", err)
	}
}

func TestSaveErrorsBeforeReplacement(t *testing.T) {
	for _, invalidVersion := range []int{0, 2} {
		t.Run("version="+strconv.Itoa(invalidVersion), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			original := []byte(`{"version":1,"streak":7}`)
			writeFixture(t, path, original)
			invalid := Default()
			invalid.Version = invalidVersion
			if err := Save(path, invalid); err == nil {
				t.Error("expected unsupported version error")
			}
			assertFileBytes(t, path, original)
			assertNoTemporaryFiles(t, filepath.Dir(path), "progress.json")
		})
	}
	t.Run("serialization", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "progress.json")
		original := []byte(`{"version":1,"streak":7}`)
		writeFixture(t, path, original)
		invalid := Default()
		when := time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)
		invalid.LastSession = &when
		if err := Save(path, invalid); err == nil {
			t.Fatal("expected JSON serialization error")
		}
		assertFileBytes(t, path, original)
		assertNoTemporaryFiles(t, filepath.Dir(path), "progress.json")
	})
	t.Run("parent is a file", func(t *testing.T) {
		parent := filepath.Join(t.TempDir(), "parent")
		writeFixture(t, parent, []byte("keep"))
		if err := Save(filepath.Join(parent, "progress.json"), Default()); err == nil {
			t.Fatal("expected path error")
		}
		assertFileBytes(t, parent, []byte("keep"))
	})
}

func TestSaveCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new", "nested", "progress.json")
	if err := Save(path, Default()); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil || !reflect.DeepEqual(got, Default()) {
		t.Fatalf("creation: %#v, %v", got, err)
	}
	assertNoTemporaryFiles(t, filepath.Dir(path), "progress.json")
}

func TestReplaceFileFailuresPreserveDestinationAndCleanTemporaryFile(t *testing.T) {
	for _, failure := range []string{"partial write", "closed file"} {
		t.Run(failure, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "progress.json")
			original := []byte(`{"version":1,"streak":7}`)
			writeFixture(t, path, original)
			injected := errors.New("write interrupted")
			err := replaceFile(path, 0o600, func(file *os.File) error {
				if filepath.Dir(file.Name()) != dir {
					t.Fatalf("temporary file is on a different directory: %s", file.Name())
				}
				if _, err := file.Write([]byte(`{"version":`)); err != nil {
					return err
				}
				if failure == "partial write" {
					return injected
				}
				if err := file.Close(); err != nil {
					return err
				}
				_, err := file.Write([]byte("1}"))
				return err
			})
			if err == nil {
				t.Fatal("expected a write failure")
			}
			if failure == "partial write" && !errors.Is(err, injected) {
				t.Fatalf("lost write error: %v", err)
			}
			if failure == "closed file" && !errors.Is(err, os.ErrClosed) {
				t.Fatalf("lost file error: %v", err)
			}
			assertFileBytes(t, path, original)
			assertNoTemporaryFiles(t, dir, "progress.json")
		})
	}
	t.Run("rename", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "progress.json")
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatal(err)
		}
		err := replaceFile(path, 0o600, func(file *os.File) error {
			_, err := file.Write([]byte(`{"version":1}`))
			return err
		})
		if err == nil {
			t.Fatal("expected rename error when replacing a directory")
		}
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			t.Fatalf("original directory changed: %v", err)
		}
		assertNoTemporaryFiles(t, dir, "progress.json")
	})
}

func TestResetRemainsExplicitlyDestructive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	writeFixture(t, path, []byte("invalid JSON"))
	if err := Reset(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Reset did not remove the selected file: %v", err)
	}
	if err := Reset(path); err != nil {
		t.Fatalf("Reset missing file: %v", err)
	}
}

func TestReplaceFileSyncAndCloseFailures(t *testing.T) {
	for _, operation := range []string{"sync", "close"} {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "progress.json")
			original := []byte(`{"version":1,"streak":7}`)
			writeFixture(t, path, original)
			injected := errors.New(operation + " failed")
			var calls []string
			err := replaceFileWithLifecycle(path, 0o600, func(file *os.File) error {
				_, err := file.Write([]byte(`{"version":1,"streak":8}`))
				return err
			}, func(file *os.File) error {
				calls = append(calls, "sync")
				assertFileBytes(t, path, original)
				if operation == "sync" {
					return injected
				}
				return file.Sync()
			}, func(file *os.File) error {
				calls = append(calls, "close")
				assertFileBytes(t, path, original)
				// A real Close may report an error after releasing the descriptor.
				return errors.Join(file.Close(), injected)
			})
			if !errors.Is(err, injected) || !strings.Contains(err.Error(), operation+" temporary file") {
				t.Fatalf("lost %s error: %v", operation, err)
			}
			wantCalls := []string{"sync"}
			if operation == "close" {
				wantCalls = append(wantCalls, "close")
			}
			if !reflect.DeepEqual(calls, wantCalls) {
				t.Fatalf("operations reached: %v, want %v", calls, wantCalls)
			}
			assertFileBytes(t, path, original)
			assertNoTemporaryFiles(t, dir, "progress.json")
		})
	}
}

func TestReplaceFileSyncsAndClosesBeforeRename(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")
	original := []byte(`{"version":1,"streak":7}`)
	replacement := []byte(`{"version":1,"streak":8}`)
	writeFixture(t, path, original)
	var calls []string
	err := replaceFileWithLifecycle(path, 0o600, func(file *os.File) error {
		_, err := file.Write(replacement)
		return err
	}, func(file *os.File) error {
		calls = append(calls, "sync")
		assertFileBytes(t, path, original)
		assertFileBytes(t, file.Name(), replacement)
		return file.Sync()
	}, func(file *os.File) error {
		calls = append(calls, "close")
		assertFileBytes(t, path, original)
		return file.Close()
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, []string{"sync", "close"}) {
		t.Fatalf("operations before rename: %v", calls)
	}
	assertFileBytes(t, path, replacement)
	assertNoTemporaryFiles(t, dir, "progress.json")
}

func writeFixture(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertFileBytes(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("original file changed: got %q, want %q (%v)", got, want, err)
	}
}

func assertNoTemporaryFiles(t *testing.T, dir, want string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != want {
		t.Fatalf("unexpected files after Save: %v", entries)
	}
}
