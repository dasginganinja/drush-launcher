package drushlauncher

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
)

type ComposerJson struct {
	Config struct {
		BinDir string `json:"bin-dir"`
	} `json:"config"`
}

func FindDrupalRoot(path string) (string, error) {
  var drushDir string

	// Check if the vendor/bin/drush directory exists in the current directory
	binDir := GetComposerBinDir(path)

	if (binDir != "") {
		drushDir := filepath.Join(path, binDir, "drush")
		if _, err := os.Stat(drushDir); err == nil {
			// Drupal root found, return the current directory
			return path, nil
		}
	}
	// Move one level up the directory tree
	parentDir := filepath.Dir(path)
	if parentDir == path {
		// If we reached the root directory, stop traversing
		return "", fmt.Errorf("Drupal root not found")
	}

	// Check if the drush exec file exists in the parent directory
	binDir = GetComposerBinDir(parentDir)
	if (binDir != "") {
		drushDir = filepath.Join(path, binDir, "drush")
		if _, err := os.Stat(drushDir); err == nil {
			// Drupal root found in the parent directory, return the parent directory
			return parentDir, nil
		}
	}

	// Recursively continue searching in the parent directory
	return FindDrupalRoot(parentDir)
}

func GetComposerBinDir(path string) string {
	composerJsonPath := filepath.Join(path, "composer.json")

	// Check if composer.json and composer.json exist
	if _, err := os.Stat(composerJsonPath); os.IsNotExist(err) {
		return ""
	}
	composerJsonFile, err := os.Open(composerJsonPath)

	// Check if we can open composer.json
	if err != nil {
		fmt.Println("Error opening composer.json:", err)
		return ""
	}

	// Close the file when we are done.
	defer composerJsonFile.Close()

	// Read the file into a byte array
	byteValue, _ := ioutil.ReadAll(composerJsonFile)

	// Parse composer.json
	var composerJsonValues ComposerJson
	json.Unmarshal(byteValue, &composerJsonValues)

	// Check if the bin-dir flag exists
	if composerJsonValues.Config.BinDir == "" {
		// If the bin-dir flag does not exist, use the default vendor/bin
		return filepath.Join("vendor", "bin")
	}

	// If the bin-dir flag exists, use the value of the flag
	binDir := filepath.Join(composerJsonValues.Config.BinDir)

	return binDir
}