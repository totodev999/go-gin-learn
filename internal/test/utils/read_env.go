package test_utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Find the directory where the file "go.mod" is placed.
func ReadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	root := FindProjectRoot(wd)
	envPath := filepath.Join(root, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: %s not found, falling back to system environment variables", envPath)
	}

}
