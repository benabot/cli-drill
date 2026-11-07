package progress

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Progress struct {
	Version            int                     `json:"version"`
	CompletedExercises []string                `json:"completed_exercises"`
	ChapterScores      map[string]ChapterScore `json:"chapter_scores"`
	LastSession        *time.Time              `json:"last_session"`
	Streak             int                     `json:"streak"`
}

type ChapterScore struct {
	Attempts int `json:"attempts"`
	Correct  int `json:"correct"`
}

func Default() Progress {
	return Progress{
		Version:            1,
		CompletedExercises: []string{},
		ChapterScores:      map[string]ChapterScore{},
		LastSession:        nil,
		Streak:             0,
	}
}

func Load(path string) (Progress, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return Progress{}, fmt.Errorf("load progress %q: %w", path, err)
	}

	initial := Default()
	progress := &initial
	if err := json.Unmarshal(data, &progress); err != nil {
		return Progress{}, fmt.Errorf("load progress %q: invalid JSON: %w", path, err)
	}
	if progress == nil {
		return Progress{}, fmt.Errorf("load progress %q: expected a JSON object, got null", path)
	}
	if progress.Version != 1 {
		return Progress{}, fmt.Errorf("load progress %q: unsupported version %d (expected 1)", path, progress.Version)
	}
	if progress.ChapterScores == nil {
		progress.ChapterScores = map[string]ChapterScore{}
	}
	return *progress, nil
}

func Save(path string, progress Progress) error {
	if progress.Version != 1 {
		return fmt.Errorf("save progress %q: unsupported version %d (expected 1)", path, progress.Version)
	}
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return fmt.Errorf("save progress %q: encode JSON: %w", path, err)
	}
	data = append(data, '\n')

	mode := os.FileMode(0o600)
	info, err := os.Lstat(path)
	if err == nil {
		// Replacing a symlink would replace the link, not its target. Refuse it
		// and other non-regular files instead of silently changing their meaning.
		if !info.Mode().IsRegular() {
			return fmt.Errorf("save progress %q: expected a regular file", path)
		}
		if _, err := Load(path); err != nil {
			return fmt.Errorf("save progress %q: refusing to overwrite unreadable or invalid progress: %w", path, err)
		}
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("save progress %q: inspect existing file: %w", path, err)
	}
	if err := replaceFile(path, mode, func(file *os.File) error {
		_, err := file.Write(data)
		return err
	}); err != nil {
		return fmt.Errorf("save progress %q: %w", path, err)
	}
	return nil
}

// replaceFile writes to a sibling temporary file, leaving the destination intact
// until the write, sync and close have all succeeded. Rename is atomic on Unix;
// it does not coordinate concurrent writers or guarantee power-loss durability.
func replaceFile(path string, mode os.FileMode, write func(*os.File) error) (result error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".cli-drill-progress-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := file.Name()
	defer func() {
		if err := file.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			result = errors.Join(result, fmt.Errorf("close temporary file: %w", err))
		}
		if temporaryPath != "" {
			if err := os.Remove(temporaryPath); err != nil {
				result = errors.Join(result, fmt.Errorf("remove temporary file: %w", err))
			}
		}
	}()
	if err := write(file); err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}
	if err := file.Chmod(mode); err != nil {
		return fmt.Errorf("set temporary file permissions: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync temporary file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	temporaryPath = ""
	return nil
}

func Reset(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (p *Progress) MarkCompleted(chapterID, itemID string, correct bool) {
	if p.Version == 0 {
		*p = Default()
	}
	if p.ChapterScores == nil {
		p.ChapterScores = map[string]ChapterScore{}
	}

	key := chapterID + "/" + itemID
	if correct && !contains(p.CompletedExercises, key) {
		p.CompletedExercises = append(p.CompletedExercises, key)
	}

	score := p.ChapterScores[chapterID]
	score.Attempts++
	if correct {
		score.Correct++
	}
	p.ChapterScores[chapterID] = score

	now := time.Now().UTC()
	p.LastSession = &now
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
