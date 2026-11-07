package chapter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/benabot/cli-drill/internal/exercise"
	"gopkg.in/yaml.v3"
)

type Chapter struct {
	ID          string `json:"id" yaml:"id"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Items       []Item `json:"items" yaml:"items"`
}

type Item struct {
	ID           string          `json:"id" yaml:"id"`
	Type         string          `json:"type" yaml:"type"`
	ExerciseType exercise.Type   `json:"exercise_type" yaml:"exercise_type"`
	Prompt       string          `json:"prompt" yaml:"prompt"`
	Choices      []string        `json:"choices,omitempty" yaml:"choices,omitempty"`
	Answer       exercise.Answer `json:"answer" yaml:"answer"`
	Explanation  string          `json:"explanation,omitempty" yaml:"explanation,omitempty"`
}

func Load(r io.Reader) (Chapter, error) {
	var chapter Chapter
	if err := yaml.NewDecoder(r).Decode(&chapter); err != nil {
		return Chapter{}, err
	}
	if err := chapter.Validate(); err != nil {
		return Chapter{}, err
	}
	return chapter, nil
}

func LoadFS(fsys fs.FS, pattern string) ([]Chapter, error) {
	paths, err := fs.Glob(fsys, pattern)
	if err != nil {
		return nil, err
	}
	var chapters []Chapter
	for _, path := range paths {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, err
		}
		chapter, err := Load(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}

func (c Chapter) Validate() error {
	if c.ID == "" {
		return errors.New("chapter id is required")
	}
	if c.Title == "" {
		return fmt.Errorf("chapter %s: title is required", c.ID)
	}
	for _, item := range c.Items {
		if item.ID == "" {
			return fmt.Errorf("chapter %s: item id is required", c.ID)
		}
		if item.Type == "" {
			return fmt.Errorf("chapter %s item %s: type is required", c.ID, item.ID)
		}
		if item.ExerciseType == "" {
			return fmt.Errorf("chapter %s item %s: exercise_type is required", c.ID, item.ID)
		}
		if item.Prompt == "" {
			return fmt.Errorf("chapter %s item %s: prompt is required", c.ID, item.ID)
		}
		if item.Answer.Primary == "" {
			return fmt.Errorf("chapter %s item %s: answer.primary is required", c.ID, item.ID)
		}
	}
	return nil
}
