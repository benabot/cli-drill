package detect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveConfiguredPath(root, configured string) (string, error) {
	root = ExpandHome(root)
	configured = ExpandHome(configured)

	if filepath.IsAbs(configured) {
		rootAbs, err := filepath.Abs(root)
		if err != nil {
			return "", err
		}
		resolved := filepath.Clean(configured)
		if !isWithin(rootAbs, resolved) {
			return "", fmt.Errorf("absolute path is outside dotfiles root: %s", configured)
		}
		return resolved, nil
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	resolved := filepath.Clean(filepath.Join(rootAbs, configured))
	if !isWithin(rootAbs, resolved) {
		return "", fmt.Errorf("configured path escapes dotfiles root: %s", configured)
	}
	return resolved, nil
}

func IsExcluded(path string, excludes []string) bool {
	path = cleanForCompare(path)
	for _, exclude := range excludes {
		if exclude == "" {
			continue
		}
		exclude = cleanForCompare(exclude)
		if path == exclude || isWithin(exclude, path) {
			return true
		}
	}
	return false
}

func ExpandHome(path string) string {
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	}
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return path
}

func cleanForCompare(path string) string {
	path = ExpandHome(path)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	return filepath.Clean(path)
}

func isWithin(root, child string) bool {
	root = filepath.Clean(root)
	child = filepath.Clean(child)
	rel, err := filepath.Rel(root, child)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
