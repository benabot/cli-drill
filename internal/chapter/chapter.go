package chapter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/benabot/cli-drill/internal/exercise"
	"gopkg.in/yaml.v3"
)

type Chapter struct {
	ID            string        `json:"id" yaml:"id"`
	Title         string        `json:"title" yaml:"title"`
	Description   string        `json:"description" yaml:"description"`
	Family        Family        `json:"family,omitempty" yaml:"family,omitempty"`
	LearningOrder LearningOrder `json:"learning_order,omitempty" yaml:"learning_order,omitempty"`
	Items         []Item        `json:"items" yaml:"items"`
}

type Item struct {
	ID           string          `json:"id" yaml:"id"`
	Type         string          `json:"type" yaml:"type"`
	ExerciseType exercise.Type   `json:"exercise_type" yaml:"exercise_type"`
	Difficulty   Difficulty      `json:"difficulty,omitempty" yaml:"difficulty,omitempty"`
	Tools        []string        `json:"tools,omitempty" yaml:"tools,omitempty"`
	Prompt       string          `json:"prompt" yaml:"prompt"`
	Choices      []string        `json:"choices,omitempty" yaml:"choices,omitempty"`
	Answer       exercise.Answer `json:"answer" yaml:"answer"`
	Explanation  string          `json:"explanation,omitempty" yaml:"explanation,omitempty"`
}

// MatchAnswer selects the correction rule from the existing pedagogical metadata.
// CLI and TUI must use this entry point so they cannot select different rules.
func (item Item) MatchAnswer(input string) bool {
	switch {
	case item.ExerciseType == exercise.TypeKeySequence:
		return exercise.MatchKeySequence(input, item.Answer)
	case item.ExerciseType == exercise.TypeSimpleShellSim:
		return exercise.MatchShellAnswer(input, item.Answer)
	case item.Type == "shortcut":
		return exercise.MatchKeySequence(input, item.Answer)
	case item.Type == "tool" || item.Type == "concept":
		return exercise.MatchConceptAnswer(input, item.Answer)
	default:
		return exercise.MatchAnswer(input, item.Answer)
	}
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
	ids := make(map[string]bool)
	for _, path := range paths {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, err
		}
		chapter, err := Load(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if ids[chapter.ID] {
			return nil, fmt.Errorf("%s: duplicate chapter id %q", path, chapter.ID)
		}
		ids[chapter.ID] = true
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
	switch c.Family {
	case "", FamilyUnixShell, FamilyKeyboardShortcuts, FamilyModernCLIMacOS:
	default:
		return fmt.Errorf("chapter %s: invalid family %q", c.ID, c.Family)
	}
	switch c.LearningOrder {
	case "", LearningOrderGuided, LearningOrderIndependent:
	default:
		return fmt.Errorf("chapter %s: invalid learning_order %q", c.ID, c.LearningOrder)
	}
	ids := make(map[string]bool)
	for _, item := range c.Items {
		if item.ID == "" {
			return fmt.Errorf("chapter %s: item id is required", c.ID)
		}
		if ids[item.ID] {
			return fmt.Errorf("chapter %s: duplicate item id %q", c.ID, item.ID)
		}
		ids[item.ID] = true
		if item.Type == "" {
			return fmt.Errorf("chapter %s item %s: type is required", c.ID, item.ID)
		}
		if item.ExerciseType == "" {
			return fmt.Errorf("chapter %s item %s: exercise_type is required", c.ID, item.ID)
		}
		switch item.ExerciseType {
		case exercise.TypeFreeAnswer, exercise.TypeMultipleChoice, exercise.TypeScenario, exercise.TypeSimpleShellSim, exercise.TypeKeySequence:
		default:
			return fmt.Errorf("chapter %s item %s: invalid exercise_type %q", c.ID, item.ID, item.ExerciseType)
		}
		switch item.Difficulty {
		case "", DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced:
		default:
			return fmt.Errorf("chapter %s item %s: invalid difficulty %q", c.ID, item.ID, item.Difficulty)
		}
		tools := make(map[string]bool)
		for _, tool := range item.Tools {
			if strings.TrimSpace(tool) == "" || tools[tool] {
				return fmt.Errorf("chapter %s item %s: tools must contain non-empty, unique identifiers (got %q)", c.ID, item.ID, tool)
			}
			tools[tool] = true
		}
		if item.Prompt == "" {
			return fmt.Errorf("chapter %s item %s: prompt is required", c.ID, item.ID)
		}
		if item.Answer.Primary == "" {
			return fmt.Errorf("chapter %s item %s: answer.primary is required", c.ID, item.ID)
		}
		if item.ExerciseType == exercise.TypeMultipleChoice {
			choices := make(map[string]bool)
			hasAnswer := false
			for _, choice := range item.Choices {
				if strings.TrimSpace(choice) == "" || choices[choice] {
					return fmt.Errorf("chapter %s item %s: choices must be non-empty and unique", c.ID, item.ID)
				}
				choices[choice] = true
				hasAnswer = hasAnswer || item.MatchAnswer(choice)
			}
			if !hasAnswer {
				return fmt.Errorf("chapter %s item %s: choices must include an accepted answer", c.ID, item.ID)
			}
		}
	}
	return nil
}
