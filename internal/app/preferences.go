package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/benabot/cli-drill/internal/chapter"
)

type preferences struct {
	Version           int                `json:"version"`
	EnabledTools      []string           `json:"enabled_tools"`
	DefaultDifficulty chapter.Difficulty `json:"default_difficulty"`
}

func defaultPreferences() preferences {
	return preferences{Version: 1, EnabledTools: []string{}, DefaultDifficulty: chapter.DifficultyBeginner}
}

func validDefaultDifficulty(level chapter.Difficulty) bool {
	switch level {
	case chapter.DifficultyBeginner, chapter.DifficultyIntermediate, chapter.DifficultyAdvanced:
		return true
	default:
		return false
	}
}

func loadPreferences(path string) (preferences, error) {
	state, err := readPreferences(path)
	if err != nil {
		return defaultPreferences(), fmt.Errorf("load preferences %q: %w", path, err)
	}
	return state, nil
}

func readPreferences(path string) (preferences, error) {
	state := defaultPreferences()
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if !info.Mode().IsRegular() {
		return state, errors.New("expected a regular file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return state, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return state, fmt.Errorf("invalid JSON: %w", err)
	}
	if fields == nil {
		return state, errors.New("expected a JSON object")
	}
	raw, present := fields["version"]
	if !present {
		return state, errors.New("missing version (expected 1)")
	}
	var version *int
	if err := json.Unmarshal(raw, &version); err != nil || version == nil {
		return state, errors.New("version must be an integer")
	}
	if *version != state.Version {
		return state, fmt.Errorf("unsupported version %d (expected 1)", *version)
	}
	for name := range fields {
		switch name {
		case "version", "enabled_tools", "default_difficulty":
		default:
			return state, fmt.Errorf("unknown field %q for version 1", name)
		}
	}
	if raw, present := fields["enabled_tools"]; present {
		var tools []*string
		if err := json.Unmarshal(raw, &tools); err != nil || tools == nil {
			return state, errors.New("enabled_tools must be an array of strings")
		}
		for _, tool := range tools {
			if tool == nil {
				return state, errors.New("enabled_tools must be an array of strings")
			}
			state.EnabledTools = append(state.EnabledTools, *tool)
		}
	}
	if raw, present := fields["default_difficulty"]; present {
		state.DefaultDifficulty = ""
		if err := json.Unmarshal(raw, &state.DefaultDifficulty); err != nil {
			return state, errors.New("default_difficulty must be beginner, intermediate or advanced")
		}
	}
	return state, state.validate()
}

func (p preferences) validate() error {
	if p.Version != 1 {
		return fmt.Errorf("unsupported version %d (expected 1)", p.Version)
	}
	if !validDefaultDifficulty(p.DefaultDifficulty) {
		return errors.New("default_difficulty must be beginner, intermediate or advanced")
	}
	for i, tool := range p.EnabledTools {
		if strings.TrimSpace(tool) == "" || slices.Contains(p.EnabledTools[:i], tool) {
			return errors.New("enabled_tools must contain non-empty, unique identifiers")
		}
	}
	return nil
}

func samePreferences(a, b preferences) bool {
	return a.Version == b.Version && a.DefaultDifficulty == b.DefaultDifficulty && slices.Equal(a.EnabledTools, b.EnabledTools)
}

// savePreferences rechecks the destination before replacing it. As with
// progress, atomic replacement assumes a single writer; it is not a lock.
func savePreferences(path string, original, state preferences) error {
	if err := writePreferences(path, original, state); err != nil {
		return fmt.Errorf("save preferences %q: %w", path, err)
	}
	return nil
}

func writePreferences(path string, original, state preferences) error {
	if err := state.validate(); err != nil {
		return err
	}
	current, err := readPreferences(path)
	if err != nil {
		return fmt.Errorf("refusing to overwrite unreadable or invalid preferences: %w", err)
	}
	if !samePreferences(current, original) {
		return errors.New("preferences changed on disk; reload before saving")
	}

	// Rename alone does not check write permission on the file it replaces.
	existing, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err == nil {
		if err := existing.Close(); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if state.EnabledTools == nil {
		state.EnabledTools = []string{}
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return replacePreferencesFile(path, append(data, '\n'))
}

func replacePreferencesFile(path string, data []byte) (result error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".cli-drill-preferences-*.tmp")
	if err != nil {
		return err
	}
	temp := file.Name()
	defer func() {
		if err := file.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			result = errors.Join(result, fmt.Errorf("close temporary file: %w", err))
		}
		if temp != "" {
			if err := os.Remove(temp); err != nil {
				result = errors.Join(result, fmt.Errorf("remove temporary file: %w", err))
			}
		}
	}()
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(temp, path); err != nil {
		return err
	}
	temp = ""
	return nil
}
