package detect

import (
	"path/filepath"
	"testing"
)

func TestIsExcludedMatchesExactPathAndChildren(t *testing.T) {
	root := t.TempDir()
	secretRoot := filepath.Join(root, "secrets")
	excludes := []string{secretRoot}

	if !IsExcluded(secretRoot, excludes) {
		t.Fatal("expected exact excluded path to match")
	}
	if !IsExcluded(filepath.Join(secretRoot, "token.txt"), excludes) {
		t.Fatal("expected child path to match")
	}
	if IsExcluded(filepath.Join(root, "notes", "token.txt"), excludes) {
		t.Fatal("did not expect unrelated path to match")
	}
}

func TestResolveConfiguredPathKeepsRelativeSourcesUnderRoot(t *testing.T) {
	root := t.TempDir()

	got, err := ResolveConfiguredPath(root, "zsh/modules/aliases.zsh")
	if err != nil {
		t.Fatalf("ResolveConfiguredPath returned error: %v", err)
	}

	want := filepath.Join(root, "zsh", "modules", "aliases.zsh")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
