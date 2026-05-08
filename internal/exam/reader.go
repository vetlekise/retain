package exam

import (
	"errors"
	"os"
	"path/filepath"
)

// ReadExam validates that the given path has a .json extension and returns its contents.
// Returns an error if the extension is invalid or the file cannot be read.
func ReadExam(path string) ([]byte, error) {
	if filepath.Ext(path) != ".json" {
		return nil, errors.New("file extension must be '.json'")
	}
	return os.ReadFile(path)
}
