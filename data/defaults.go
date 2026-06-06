package data

import (
	"embed"
	"io/fs"
)

//go:embed chapters/*.yaml
var defaultChapters embed.FS

func Chapters() fs.FS {
	chapters, err := fs.Sub(defaultChapters, "chapters")
	if err != nil {
		panic(err)
	}
	return chapters
}
