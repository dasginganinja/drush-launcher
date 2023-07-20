package drushlauncher

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindDrupalRoot(path string) (string, error) {
	// Start traversing from the provided path to the root directory
	currentDir := path
	for {
		// Check if the vendor/bin/drush directory exists in the current directory
		drushDir := filepath.Join(currentDir, "vendor", "bin", "drush")
		if _, err := os.Stat(drushDir); err == nil {
			// Drupal root found, traverse one level up to get the root directory
			drupalRoot := filepath.Dir(currentDir)
			return drupalRoot, nil
		}

		// Move one level up the directory tree
		parentDir := filepath.Dir(currentDir)
		// If we reached the root directory, stop traversing
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	// Drupal root not found in the entire directory tree
	return "", fmt.Errorf("Drupal root not found")
}
