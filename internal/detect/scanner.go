package detect

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/benabot/cli-drill/internal/catalog"
	"github.com/benabot/cli-drill/internal/config"
	"github.com/benabot/cli-drill/internal/markdown"
	"github.com/benabot/cli-drill/internal/shell/zsh"
)

type Warning struct {
	Path    string
	Message string
}

func Scan(cfg config.Config) (catalog.Catalog, []Warning) {
	c := catalog.New()
	var warnings []Warning
	if err := cfg.Validate(); err != nil {
		return c, []Warning{{Message: err.Error()}}
	}

	for _, source := range cfg.Paths.Aliases {
		path, ok := safeSourcePath(cfg, source, &warnings)
		if !ok {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			warnings = append(warnings, Warning{Path: path, Message: err.Error()})
			continue
		}
		aliases, err := zsh.ParseAliases(bytes.NewReader(data), source)
		if err != nil {
			warnings = append(warnings, Warning{Path: path, Message: err.Error()})
			continue
		}
		for _, entry := range zsh.AliasEntries(aliases) {
			c.Add(entry)
			addToolFromCommand(&c, entry.Command, entry.Source)
		}
	}

	for _, source := range cfg.Paths.Functions {
		path, ok := safeSourcePath(cfg, source, &warnings)
		if !ok {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			warnings = append(warnings, Warning{Path: path, Message: err.Error()})
			continue
		}
		functions, err := zsh.ParseFunctions(bytes.NewReader(data), source)
		if err != nil {
			warnings = append(warnings, Warning{Path: path, Message: err.Error()})
			continue
		}
		for _, entry := range zsh.FunctionEntries(functions) {
			c.Add(entry)
		}
	}

	for _, source := range cfg.Paths.Docs {
		path, ok := safeSourcePath(cfg, source, &warnings)
		if !ok {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			warnings = append(warnings, Warning{Path: path, Message: err.Error()})
			continue
		}
		doc, err := markdown.Parse(bytes.NewReader(data))
		if err != nil {
			warnings = append(warnings, Warning{Path: path, Message: err.Error()})
			continue
		}
		for _, heading := range doc.Headings {
			c.Add(catalog.Entry{
				ID:      catalog.NormalizeID(heading.Text),
				Name:    heading.Text,
				Type:    catalog.EntryConcept,
				Summary: "Markdown heading",
				Source:  catalog.Source{Path: source, Line: heading.Line},
			})
		}
		addKnownToolsFromText(&c, string(data), source)
	}

	return c, warnings
}

func safeSourcePath(cfg config.Config, source string, warnings *[]Warning) (string, bool) {
	path, err := ResolveConfiguredPath(cfg.DotfilesPath, source)
	if err != nil {
		*warnings = append(*warnings, Warning{Path: source, Message: err.Error()})
		return "", false
	}
	if IsExcluded(path, cfg.Security.Exclude) {
		*warnings = append(*warnings, Warning{Path: path, Message: "excluded by security policy"})
		return "", false
	}
	return path, true
}

func addToolFromCommand(c *catalog.Catalog, command string, source catalog.Source) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return
	}
	name := strings.TrimPrefix(filepath.Base(fields[0]), "command")
	if name == "" || strings.ContainsAny(name, "$/") {
		return
	}
	if !isKnownTool(name) {
		return
	}
	c.Add(catalog.Entry{
		ID:      catalog.NormalizeID("tool-" + name),
		Name:    name,
		Type:    catalog.EntryTool,
		Summary: "Detected from alias command",
		Source:  source,
	})
}

func addKnownToolsFromText(c *catalog.Catalog, text, source string) {
	for _, tool := range knownTools {
		re := regexp.MustCompile(`(?i)(^|[^a-z0-9_-])` + regexp.QuoteMeta(tool) + `([^a-z0-9_-]|$)`)
		if re.FindStringIndex(text) == nil {
			continue
		}
		c.Add(catalog.Entry{
			ID:      catalog.NormalizeID("tool-" + tool),
			Name:    tool,
			Type:    catalog.EntryTool,
			Summary: "Detected in documentation",
			Source:  catalog.Source{Path: source},
		})
	}
}

func isKnownTool(name string) bool {
	name = strings.ToLower(name)
	for _, tool := range knownTools {
		if name == strings.ToLower(tool) {
			return true
		}
	}
	return false
}

var knownTools = []string{
	"atuin", "yazi", "zoxide", "fzf", "fd", "rg", "bat", "glow", "eza",
	"micro", "lazygit", "lazydocker", "delta", "sd", "tokei", "xh",
	"doggo", "gum", "viddy", "gitleaks", "hadolint", "vale",
	"markdownlint-cli2", "swiftlint", "swiftformat", "xcodegen",
}
