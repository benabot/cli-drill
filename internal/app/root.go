package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/benabot/cli-drill/internal/catalog"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/config"
	"github.com/benabot/cli-drill/internal/detect"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
	"github.com/benabot/cli-drill/internal/tui"
	"github.com/benabot/cli-drill/internal/xdg"
	"github.com/spf13/cobra"
)

type Options struct {
	ConfigPath   string
	ProgressPath string
	In           io.Reader
	Out          io.Writer
	Err          io.Writer
	DefaultFS    fs.FS
}

func NewRootCommand(defaultFS fs.FS) *cobra.Command {
	opts := Options{
		In:        os.Stdin,
		Out:       os.Stdout,
		Err:       os.Stderr,
		DefaultFS: defaultFS,
	}

	root := &cobra.Command{
		Use:   "cli-drill",
		Short: "Train with your shell, dotfiles and CLI workflows",
		Long:  "cli-drill scans dotfiles statically, builds a typed directory and trains from editable YAML chapters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, opts)
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			entries := chapter.ToCatalog(chapters).Entries()
			return tui.Run(chapters, entries)
		},
	}
	root.PersistentFlags().StringVar(&opts.ConfigPath, "config", "", "path to config TOML")
	root.PersistentFlags().StringVar(&opts.ProgressPath, "progress", "", "path to progress JSON")

	root.AddCommand(newInitCommand(&opts))
	root.AddCommand(newScanCommand(&opts))
	root.AddCommand(newGenerateCommand(&opts))
	root.AddCommand(newChaptersCommand(&opts))
	root.AddCommand(newTrainCommand(&opts))
	root.AddCommand(newDirectoryCommand(&opts))
	root.AddCommand(newSearchCommand(&opts))
	root.AddCommand(newShowCommand(&opts))
	root.AddCommand(newStatsCommand(&opts))
	root.AddCommand(newResetCommand(&opts))

	return root
}

func newInitCommand(opts *Options) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a starter config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			path := configPath(runOpts)
			if err := config.Save(path, config.Default(), force); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(runOpts.Out, "Created config: %s\n", path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing config")
	return cmd
}

func newScanCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "Scan configured dotfiles statically",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			cfg, err := loadConfig(runOpts)
			if err != nil {
				return err
			}
			entries, warnings := detect.Scan(cfg)
			printWarnings(runOpts.Err, warnings)
			printEntries(runOpts.Out, entries.Entries())
			return nil
		},
	}
}

func newGenerateCommand(opts *Options) *cobra.Command {
	var outDir string
	var force bool
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate editable YAML chapters from the current scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			cfg, err := loadConfig(runOpts)
			if err != nil {
				return err
			}
			entries, warnings := detect.Scan(cfg)
			printWarnings(runOpts.Err, warnings)
			chapters := chapter.GenerateFromCatalog(entries.Entries())
			if outDir == "" {
				outDir = xdg.ChaptersDir()
			}
			if err := chapter.WriteDir(outDir, chapters, force); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(runOpts.Out, "Wrote %d chapters to %s\n", len(chapters), outDir)
			return nil
		},
	}
	cmd.Flags().StringVar(&outDir, "out", "", "output directory for generated chapters")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing chapters")
	return cmd
}

func newChaptersCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "chapters",
		Short: "List available chapters",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			for _, chapter := range chapters {
				_, _ = fmt.Fprintf(runOpts.Out, "%s\t%s\t%d exercises\n", chapter.ID, chapter.Title, len(chapter.Items))
			}
			return nil
		},
	}
}

func newTrainCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "train [chapter-id]",
		Short: "Train from a chapter in the terminal",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			if len(chapters) == 0 {
				return errors.New("no chapters available")
			}
			selected := chapters[0]
			if len(args) == 1 {
				var found bool
				for _, candidate := range chapters {
					if candidate.ID == args[0] {
						selected = candidate
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("chapter not found: %s", args[0])
				}
			}
			return runTraining(runOpts, selected)
		},
	}
}

func newDirectoryCommand(opts *Options) *cobra.Command {
	var entryType string
	cmd := &cobra.Command{
		Use:   "directory",
		Short: "List the typed directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			entries, err := catalogForDirectory(runOpts)
			if err != nil {
				return err
			}
			if entryType != "" {
				printEntries(runOpts.Out, entries.FilterByType(catalog.EntryType(entryType)))
				return nil
			}
			printEntries(runOpts.Out, entries.Entries())
			return nil
		},
	}
	cmd.Flags().StringVar(&entryType, "type", "", "filter by entry type")
	return cmd
}

func newSearchCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search the typed directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			entries, err := catalogForDirectory(runOpts)
			if err != nil {
				return err
			}
			printEntries(runOpts.Out, entries.Search(args[0]))
			return nil
		},
	}
}

func newShowCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show <entry>",
		Short: "Show a directory entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			entries, err := catalogForDirectory(runOpts)
			if err != nil {
				return err
			}
			entry, ok := entries.Find(args[0])
			if !ok {
				return fmt.Errorf("entry not found: %s", args[0])
			}
			_, _ = fmt.Fprintf(runOpts.Out, "ID: %s\nName: %s\nType: %s\nSummary: %s\n", entry.ID, entry.Name, entry.Type, entry.Summary)
			if entry.Command != "" {
				_, _ = fmt.Fprintf(runOpts.Out, "Command: %s\n", entry.Command)
			}
			return nil
		},
	}
}

func newStatsCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show local progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			state, err := progress.Load(progressPath(runOpts))
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(runOpts.Out, "Completed exercises: %d\n", len(state.CompletedExercises))
			_, _ = fmt.Fprintf(runOpts.Out, "Chapters with attempts: %d\n", len(state.ChapterScores))
			_, _ = fmt.Fprintf(runOpts.Out, "Streak: %d\n", state.Streak)
			return nil
		},
	}
}

func newResetCommand(opts *Options) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset local progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			if !yes {
				return errors.New("reset requires --yes")
			}
			if err := progress.Reset(progressPath(runOpts)); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(runOpts.Out, "Progress reset.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm progress reset")
	return cmd
}

func runTraining(opts Options, selected chapter.Chapter) error {
	if len(selected.Items) == 0 {
		return fmt.Errorf("chapter %s has no exercises", selected.ID)
	}

	state, err := progress.Load(progressPath(opts))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(opts.In)

	_, _ = fmt.Fprintf(opts.Out, "Chapitre: %s\n", selected.Title)
	for i, item := range selected.Items {
		_, _ = fmt.Fprintf(opts.Out, "\n%d/%d %s\n> ", i+1, len(selected.Items), item.Prompt)
		input, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		correct := exercise.MatchAnswer(input, item.Answer)
		if correct {
			_, _ = fmt.Fprintln(opts.Out, "Correct.")
		} else {
			_, _ = fmt.Fprintf(opts.Out, "Pas encore. Reponse attendue: %s\n", item.Answer.Primary)
		}
		if item.Explanation != "" {
			_, _ = fmt.Fprintln(opts.Out, item.Explanation)
		}
		state.MarkCompleted(selected.ID, item.ID, correct)
		if errors.Is(err, io.EOF) {
			break
		}
	}
	return progress.Save(progressPath(opts), state)
}

func loadConfig(opts Options) (config.Config, error) {
	path := configPath(opts)
	cfg, err := config.Load(path)
	if err == nil {
		return cfg, cfg.Validate()
	}
	if errors.Is(err, os.ErrNotExist) {
		cfg := config.Default()
		return cfg, cfg.Validate()
	}
	return config.Config{}, err
}

func loadChapters(opts Options) ([]chapter.Chapter, error) {
	userDir := xdg.ChaptersDir()
	if entries, err := filepath.Glob(filepath.Join(userDir, "*.yaml")); err == nil && len(entries) > 0 {
		return chapter.LoadDir(userDir)
	}
	return chapter.LoadFS(opts.DefaultFS, "*.yaml")
}

func catalogForDirectory(opts Options) (catalog.Catalog, error) {
	if shouldUseScanCatalog(opts) {
		cfg, err := loadConfig(opts)
		if err != nil {
			return catalog.Catalog{}, err
		}
		entries, _ := detect.Scan(cfg)
		return entries, nil
	}
	chapters, err := loadChapters(opts)
	if err != nil {
		return catalog.Catalog{}, err
	}
	return chapter.ToCatalog(chapters), nil
}

func commandOptions(cmd *cobra.Command, opts Options) Options {
	opts.In = cmd.InOrStdin()
	opts.Out = cmd.OutOrStdout()
	opts.Err = cmd.ErrOrStderr()
	return opts
}

func shouldUseScanCatalog(opts Options) bool {
	if opts.ConfigPath != "" {
		return true
	}
	if _, err := os.Stat(xdg.ConfigFile()); err == nil {
		return true
	}
	return false
}

func configPath(opts Options) string {
	if opts.ConfigPath != "" {
		return opts.ConfigPath
	}
	return xdg.ConfigFile()
}

func progressPath(opts Options) string {
	if opts.ProgressPath != "" {
		return opts.ProgressPath
	}
	return xdg.ProgressFile()
}

func printWarnings(w io.Writer, warnings []detect.Warning) {
	for _, warning := range warnings {
		if warning.Path == "" {
			_, _ = fmt.Fprintf(w, "warning: %s\n", warning.Message)
			continue
		}
		_, _ = fmt.Fprintf(w, "warning: %s: %s\n", warning.Path, warning.Message)
	}
}

func printEntries(w io.Writer, entries []catalog.Entry) {
	if len(entries) == 0 {
		_, _ = fmt.Fprintln(w, "No entries found.")
		return
	}
	for _, entry := range entries {
		summary := strings.TrimSpace(entry.Summary)
		if entry.Command != "" {
			summary = entry.Command
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", entry.Name, entry.Type, entry.ID, summary)
	}
}
