package drushlauncher

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindDrushExecutable(path string) (string, error) {
	// Check if the vendor/bin/drush directory exists in the current directory
	drushDir := filepath.Join(path, "vendor", "bin", "drush")
	if _, err := os.Stat(drushDir); err == nil {
		// Drupal root found, return the current directory
		return path, nil
	}

	// Move one level up the directory tree
	parentDir := filepath.Dir(path)
	if parentDir == path {
		// If we reached the root directory, stop traversing
		return "", fmt.Errorf("Drupal root not found")
	}

	// Check if the vendor/bin/drush directory exists in the parent directory
	drushDir = filepath.Join(parentDir, "vendor", "bin", "drush")
	if _, err := os.Stat(drushDir); err == nil {
		// Drupal root found in the parent directory, return the parent directory
		return parentDir, nil
	}

	// Recursively continue searching in the parent directory
	return FindDrushExecutable(parentDir)
}
