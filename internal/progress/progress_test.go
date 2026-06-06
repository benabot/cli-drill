package progress

import (
	"path/filepath"
	"testing"
)

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
