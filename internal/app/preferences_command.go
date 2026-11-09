package app

import (
	"fmt"
	"slices"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/xdg"
	"github.com/spf13/cobra"
)

func newPreferencesCommand(opts *Options) *cobra.Command {
	var enable, disable, level string
	cmd := &cobra.Command{
		Use:           "prefs",
		Short:         "Show or change learning preferences",
		Long:          "Show or save preferences for future sessions. Modules come from the built-in exercises and do not need to be installed. Current sessions are unchanged.",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			chapters, err := loadChapters(*opts)
			if err != nil {
				return err
			}
			available := chapter.AvailableTools(chapters)
			path := xdg.PreferencesFile()
			initial, loadErr := loadPreferences(path)
			out := cmd.OutOrStdout()
			changing := cmd.Flags().Changed("enable") || cmd.Flags().Changed("disable") || cmd.Flags().Changed("difficulty")
			if loadErr != nil {
				if changing {
					return loadErr
				}
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Preferences unavailable; showing defaults: %v\n", loadErr)
			}
			for _, tool := range chapter.UnavailableTools(available, initial.EnabledTools) {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Module %q is no longer in the corpus; saved choice kept.\n", tool)
			}
			if !changing {
				_, _ = fmt.Fprintf(out, "Default difficulty: %s\n", initial.DefaultDifficulty)
				_, _ = fmt.Fprintln(out, "CLI modules (no installation required):")
				for _, tool := range available {
					status := "disabled"
					if slices.Contains(initial.EnabledTools, tool) {
						status = "enabled"
					}
					_, _ = fmt.Fprintf(out, "  %s: %s\n", tool, status)
				}
				return nil
			}

			updated := initial
			switch {
			case cmd.Flags().Changed("enable"):
				if !slices.Contains(available, enable) {
					return fmt.Errorf("unknown module %q; run cli-drill prefs to list available modules", enable)
				}
				if !slices.Contains(initial.EnabledTools, enable) {
					updated.EnabledTools = append(slices.Clone(initial.EnabledTools), enable)
				}
			case cmd.Flags().Changed("disable"):
				if !slices.Contains(available, disable) {
					return fmt.Errorf("unknown module %q; run cli-drill prefs to list available modules", disable)
				}
				updated.EnabledTools = slices.DeleteFunc(slices.Clone(initial.EnabledTools), func(tool string) bool { return tool == disable })
			case cmd.Flags().Changed("difficulty"):
				updated.DefaultDifficulty = chapter.Difficulty(level)
				if !validDefaultDifficulty(updated.DefaultDifficulty) {
					return fmt.Errorf("invalid difficulty %q; choose beginner, intermediate or advanced", level)
				}
			}
			if samePreferences(initial, updated) {
				_, _ = fmt.Fprintln(out, "Preferences unchanged.")
				return nil
			}
			if err := savePreferences(path, initial, updated); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(out, "Preferences saved.")
			return nil
		},
	}
	cmd.Flags().StringVar(&enable, "enable", "", "enable a CLI module")
	cmd.Flags().StringVar(&disable, "disable", "", "disable a CLI module")
	cmd.Flags().StringVar(&level, "difficulty", "", "set the default difficulty (beginner, intermediate, advanced)")
	cmd.MarkFlagsMutuallyExclusive("enable", "disable", "difficulty")
	return cmd
}
