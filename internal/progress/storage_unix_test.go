//go:build darwin || linux

package progress

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveReplacesFileWithoutTruncatingPreviousInode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	original := []byte(`{"version":1,"streak":7}`)
	writeFixture(t, path, original)
	reader, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if err := Save(path, Default()); err != nil {
		t.Fatal(err)
	}
	oldBytes, err := io.ReadAll(reader)
	if err != nil || string(oldBytes) != string(original) {
		t.Fatalf("old inode was truncated: %q (%v)", oldBytes, err)
	}
	oldInfo, err := reader.Stat()
	if err != nil {
		t.Fatal(err)
	}
	newInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(oldInfo, newInfo) {
		t.Fatal("Save did not atomically replace the file")
	}
}

func TestSaveRefusesSymlinkWithoutChangingItsTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.json")
	path := filepath.Join(dir, "progress.json")
	original := []byte(`{"version":1,"streak":7}`)
	writeFixture(t, target, original)
	if err := os.Symlink("target.json", path); err != nil {
		t.Fatal(err)
	}
	if err := Save(path, Default()); err == nil {
		t.Fatal("expected explicit error for symlink destination")
	}
	link, err := os.Readlink(path)
	if err != nil || link != "target.json" {
		t.Fatalf("link was replaced: %q (%v)", link, err)
	}
	assertFileBytes(t, target, original)
}

func TestPermissionErrorsPreserveData(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root bypasses Unix permission checks")
	}
	for _, mode := range []os.FileMode{0o400, 0o444} {
		t.Run("read only "+mode.String(), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			original := []byte(`{"version":1,"streak":7}`)
			writeFixture(t, path, original)
			if err := os.Chmod(path, mode); err != nil {
				t.Fatal(err)
			}
			if _, err := Load(path); err != nil {
				t.Fatalf("read-only file must remain readable: %v", err)
			}
			if err := Save(path, Default()); !errors.Is(err, os.ErrPermission) {
				t.Errorf("Save must respect write permission: %v", err)
			}
			assertFileBytes(t, path, original)
			assertNoTemporaryFiles(t, filepath.Dir(path), "progress.json")
		})
	}
	t.Run("read", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "progress.json")
		original := []byte(`{"version":1,"streak":7}`)
		writeFixture(t, path, original)
		if err := os.Chmod(path, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
		if _, err := Load(path); !errors.Is(err, os.ErrPermission) {
			t.Fatalf("expected read permission error: %v", err)
		}
		if err := Save(path, Default()); !errors.Is(err, os.ErrPermission) {
			t.Fatalf("Save must refuse an unreadable file: %v", err)
		}
		if err := os.Chmod(path, 0o600); err != nil {
			t.Fatal(err)
		}
		assertFileBytes(t, path, original)
		assertNoTemporaryFiles(t, filepath.Dir(path), "progress.json")
	})
	t.Run("create temporary file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "progress.json")
		original := []byte(`{"version":1,"streak":7}`)
		writeFixture(t, path, original)
		if err := os.Chmod(dir, 0o500); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
		if err := Save(path, Default()); !errors.Is(err, os.ErrPermission) {
			t.Fatalf("expected write permission error: %v", err)
		}
		assertFileBytes(t, path, original)
		assertNoTemporaryFiles(t, dir, "progress.json")
	})
}

func TestSavePermissions(t *testing.T) {
	t.Run("new file is private", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "progress.json")
		if err := Save(path, Default()); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Fatalf("new file permissions = %o, want 600", info.Mode().Perm())
		}
	})
	for _, mode := range []os.FileMode{0o600, 0o640, 0o644} {
		t.Run(mode.String(), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			writeFixture(t, path, []byte(`{"version":1}`))
			if err := os.Chmod(path, mode); err != nil {
				t.Fatal(err)
			}
			if err := Save(path, Default()); err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != mode {
				t.Fatalf("permissions changed: %o, want %o", info.Mode().Perm(), mode)
			}
		})
	}
}

func TestCleanupFailureIsReported(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root bypasses Unix permission checks")
	}
	for _, operation := range []string{"write", "sync", "close"} {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "progress.json")
			original := []byte(`{"version":1,"streak":7}`)
			writeFixture(t, path, original)
			t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
			interrupted := errors.New(operation + " interrupted")
			fail := func() error {
				if err := os.Chmod(dir, 0o500); err != nil {
					t.Fatal(err)
				}
				return interrupted
			}
			err := replaceFileWithLifecycle(path, 0o600, func(file *os.File) error {
				if _, err := file.Write([]byte("partial")); err != nil {
					return err
				}
				if operation == "write" {
					return fail()
				}
				return nil
			}, func(file *os.File) error {
				if operation == "sync" {
					return fail()
				}
				return file.Sync()
			}, func(file *os.File) error {
				if err := file.Close(); err != nil {
					return err
				}
				return fail()
			})
			if !errors.Is(err, interrupted) || !errors.Is(err, os.ErrPermission) ||
				!strings.Contains(err.Error(), "remove temporary file") {
				t.Fatalf("must report both %s and cleanup errors: %v", operation, err)
			}
			assertFileBytes(t, path, original)
		})
	}
}
