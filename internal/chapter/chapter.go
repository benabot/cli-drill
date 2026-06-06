package chapter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/benabot/cli-drill/internal/catalog"
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
	ID           string            `json:"id" yaml:"id"`
	Type         catalog.EntryType `json:"type" yaml:"type"`
	ExerciseType exercise.Type     `json:"exercise_type" yaml:"exercise_type"`
	Prompt       string            `json:"prompt" yaml:"prompt"`
	Choices      []string          `json:"choices,omitempty" yaml:"choices,omitempty"`
	Answer       exercise.Answer   `json:"answer" yaml:"answer"`
	Explanation  string            `json:"explanation,omitempty" yaml:"explanation,omitempty"`
	Source       catalog.Source    `json:"source,omitempty" yaml:"source,omitempty"`
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

func LoadDir(dir string) ([]Chapter, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	var chapters []Chapter
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		chapter, loadErr := Load(file)
		closeErr := file.Close()
		if loadErr != nil {
			return nil, fmt.Errorf("%s: %w", path, loadErr)
		}
		if closeErr != nil {
			return nil, closeErr
		}
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}

func WriteDir(dir string, chapters []Chapter, force bool) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, chapter := range chapters {
		if err := chapter.Validate(); err != nil {
			return err
		}
		path := filepath.Join(dir, chapter.ID+".yaml")
		if !force {
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("chapter already exists: %s", path)
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
		data, err := yaml.Marshal(chapter)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}
	return nil
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

func ToCatalog(chapters []Chapter) catalog.Catalog {
	c := catalog.New()
	for _, chapter := range chapters {
		c.Add(catalog.Entry{
			ID:      chapter.ID,
			Name:    chapter.Title,
			Type:    catalog.EntryChapter,
			Summary: chapter.Description,
		})
		for _, item := range chapter.Items {
			c.Add(catalog.Entry{
				ID:      item.ID,
				Name:    item.Answer.Primary,
				Type:    item.Type,
				Summary: item.Prompt,
				Source:  item.Source,
				Tags:    []string{chapter.ID, string(item.ExerciseType)},
			})
		}
	}
	return c
}

func GenerateFromCatalog(entries []catalog.Entry) []Chapter {
	chapters := seedChapters()
	index := map[string]int{}
	for i := range chapters {
		index[chapters[i].ID] = i
	}

	for _, entry := range entries {
		chapterID := chapterIDForEntry(entry)
		i, ok := index[chapterID]
		if !ok {
			continue
		}
		chapters[i].Items = append(chapters[i].Items, itemForEntry(entry))
	}

	return chapters
}

func seedChapters() []Chapter {
	return []Chapter{
		{ID: "01-raccourcis-terminal", Title: "Raccourcis terminal", Description: "Memoriser les raccourcis de navigation et edition de ligne."},
		{ID: "02-navigation-shell", Title: "Navigation shell", Description: "S'entrainer aux outils de navigation du shell."},
		{ID: "03-alias-zsh", Title: "Alias ZSH", Description: "Retenir les alias ZSH utiles."},
		{ID: "04-fonctions-zsh", Title: "Fonctions ZSH", Description: "Identifier les fonctions ZSH disponibles."},
		{ID: "05-outils-quotidiens", Title: "Outils quotidiens", Description: "Reviser les outils CLI courants."},
		{ID: "06-recherche-fichiers-contenu", Title: "Recherche fichiers et contenu", Description: "S'entrainer avec fd, rg et fzf."},
		{ID: "07-lecture-preview", Title: "Lecture et preview", Description: "Utiliser les outils de lecture et de preview."},
		{ID: "08-micro", Title: "Micro", Description: "Retenir les commandes et raccourcis Micro."},
		{ID: "09-markdown-glow", Title: "Markdown et Glow", Description: "Travailler avec Markdown et Glow."},
		{ID: "10-workflows-dotfiles", Title: "Workflows dotfiles", Description: "Relier commandes, outils et workflows dotfiles."},
	}
}

func chapterIDForEntry(entry catalog.Entry) string {
	switch entry.Type {
	case catalog.EntryAlias:
		return "03-alias-zsh"
	case catalog.EntryFunction:
		return "04-fonctions-zsh"
	case catalog.EntryShortcut, catalog.EntryBinding:
		return "01-raccourcis-terminal"
	case catalog.EntryTool:
		if isSearchTool(entry.Name) {
			return "06-recherche-fichiers-contenu"
		}
		if strings.EqualFold(entry.Name, "micro") {
			return "08-micro"
		}
		if strings.EqualFold(entry.Name, "glow") {
			return "09-markdown-glow"
		}
		return "05-outils-quotidiens"
	case catalog.EntryWorkflow:
		return "10-workflows-dotfiles"
	default:
		return "10-workflows-dotfiles"
	}
}

func itemForEntry(entry catalog.Entry) Item {
	item := Item{
		ID:           entry.ID,
		Type:         entry.Type,
		ExerciseType: exercise.TypeFreeAnswer,
		Answer:       exercise.Answer{Primary: entry.Name},
		Source:       entry.Source,
	}

	switch entry.Type {
	case catalog.EntryAlias:
		item.Prompt = "Quel alias correspond a: " + entry.Command + " ?"
		item.Explanation = entry.Name + " = " + entry.Command
	case catalog.EntryFunction:
		item.Prompt = "Quelle fonction ZSH correspond a cette entree ?"
		item.Explanation = entry.Name + " est une fonction ZSH detectee statiquement."
	case catalog.EntryTool:
		item.Prompt = "Quel outil correspond a cet usage: " + entry.Summary + " ?"
		item.Explanation = entry.Name + " est un outil detecte dans les sources."
	default:
		item.Prompt = entry.Summary
		item.Explanation = entry.Summary
	}
	return item
}

func isSearchTool(name string) bool {
	switch strings.ToLower(name) {
	case "fd", "rg", "fzf", "atuin":
		return true
	default:
		return false
	}
}
