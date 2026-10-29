// Package web holds the dashboard, embedded into the grabb binary.
package web

import (
	"embed"
	"io/fs"
)

//go:embed app
var files embed.FS

// FS returns the dashboard files with app/ as the root.
func FS() fs.FS {
	sub, err := fs.Sub(files, "app")
	if err != nil {
		panic(err) // the directory is embedded at build time
	}
	return sub
}
