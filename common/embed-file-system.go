package common

import (
	"embed"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
)

// embedFileSystem wraps embed.FS for static file serving
type embedFileSystem struct {
	http.FileSystem
}

func (e *embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func (e *embedFileSystem) Open(name string) (http.File, error) {
	if name == "/" {
		// This will make sure the index page goes to NoRouter handler
		return nil, os.ErrNotExist
	}
	return e.FileSystem.Open(name)
}

// EmbedFolder creates a static.ServeFileSystem from an embedded FS
func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return &embedFileSystem{
		FileSystem: http.FS(efs),
	}
}
