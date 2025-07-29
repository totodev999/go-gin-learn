package test

import (
	"os"
	"path/filepath"
)

// Find the directory where the file "go.mod" is placed.
func FindProjectRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "." // fallback
}
