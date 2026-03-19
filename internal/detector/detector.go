package detector

import (
	"os"
	"path/filepath"
)

func DetectProject(path string) string {
	// Laravel detection
	if fileExists(filepath.Join(path, "artisan")) &&
		fileExists(filepath.Join(path, "composer.json")) {
		return "laravel"
	}

	// Node.js detection
	if fileExists(filepath.Join(path, "package.json")) {
		return "node"
	}

	// Python detection
	if fileExists(filepath.Join(path, "requirements.txt")) ||
		fileExists(filepath.Join(path, "pyproject.toml")) {
		return "python"
	}

	// Go detection
	if fileExists(filepath.Join(path, "go.mod")) {
		return "go"
	}

	return "unknown"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
