package app

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Commands now read preferences during training and statistics, even with --progress.
	dir, err := os.MkdirTemp("", "cli-drill-test-data-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.Setenv("XDG_DATA_HOME", dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		_ = os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	if err := os.RemoveAll(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		code = 1
	}
	os.Exit(code)
}
