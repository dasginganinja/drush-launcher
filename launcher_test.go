package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDrupalRoot(t *testing.T) {
	// Set up a temporary directory for the tests
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	// Test case 1: Drupal root found in a subdirectory
	testDir1 := filepath.Join(tmpDir, "testDir1", "vendor", "bin", "drush")
	os.MkdirAll(filepath.Dir(testDir1), os.ModePerm)
	os.Create(testDir1)
	drupalRoot, err := findDrupalRoot(testDir1)
	if err != nil || drupalRoot != filepath.Join(tmpDir, "testDir1") {
		t.Errorf("Test case 1 failed: Expected Drupal root %s, got %s, error: %v", filepath.Join(tmpDir, "testDir1"), drupalRoot, err)
	}

	// Test case 2: Drupal root found in the current directory
	testDir2 := filepath.Join(tmpDir, "testDir2", "vendor", "bin", "drush")
	os.MkdirAll(filepath.Dir(testDir2), os.ModePerm)
	os.Create(testDir2)
	drupalRoot, err = findDrupalRoot(testDir2)
	if err != nil || drupalRoot != filepath.Join(tmpDir, "testDir2") {
		t.Errorf("Test case 2 failed: Expected Drupal root %s, got %s, error: %v", filepath.Join(tmpDir, "testDir2"), drupalRoot, err)
	}

	// Test case 3: Drupal root not found
	testDir3 := filepath.Join(tmpDir, "testDir3", "some", "random", "path")
	drupalRoot, err = findDrupalRoot(testDir3)
	if err == nil || drupalRoot != "" {
		t.Errorf("Test case 3 failed: Expected no Drupal root, got %s, error: %v", drupalRoot, err)
	}
}
