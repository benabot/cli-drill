package progress

import (
	"encoding/json"
	"errors"
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
		return Progress{}, err
	}

	progress := Default()
	if err := json.Unmarshal(data, &progress); err != nil {
		return Progress{}, err
	}
	if progress.ChapterScores == nil {
		progress.ChapterScores = map[string]ChapterScore{}
	}
	return progress, nil
}

func Save(path string, progress Progress) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
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
