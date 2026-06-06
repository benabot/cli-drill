package main

import (
	"fmt"
	"os"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/app"
)

func main() {
	root := app.NewRootCommand(data.Chapters())
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
