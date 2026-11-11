package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/app"
)

func main() {
	root := app.NewRootCommand(data.Chapters())
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var exitStatus interface{ ExitCode() int }
		if errors.As(err, &exitStatus) {
			os.Exit(exitStatus.ExitCode())
		}
		os.Exit(1)
	}
}
