package lib

import (
	"io"
	"os"
)

// FileSystem interface for file system operations
type FileSystem interface {
	Create(path string) (io.WriteCloser, error)
}

type fs struct{}

// A thin wrapper around os package to allow mocking
// for testing purposes and to abstract file system operations
// for better testability
func NewFileSystem() FileSystem {
	return &fs{}
}

// Create creates a file at the specified path
func (f *fs) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
