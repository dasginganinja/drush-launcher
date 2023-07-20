package drush_launcher

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const drushExecutablePath = "vendor/bin/drush"

func main() {
	// Step 1: Parse command-line flags
	var rootDirFlag string
	flag.StringVar(&rootDirFlag, "root", "", "Alternative Drupal root directory")
	flag.StringVar(&rootDirFlag, "r", "", "Alternative Drupal root directory (shorthand)")
	flag.Parse()

	// Step 2: Determine the Drupal root directory
	drupalRoot := determineDrupalRoot(rootDirFlag)

	// Step 3: Execute the drush command with arguments
	cmd := exec.Command(determineDrushExecutable(drupalRoot), flag.Args()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func determineDrupalRoot(alternativeRoot string) string {
	if alternativeRoot != "" {
		// If an alternative root is provided, use it directly
		return alternativeRoot
	}

	// If no alternative root is provided, find the Drupal root as before
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Traverse upwards from the current directory to the root
	for dir != "/" && dir != "." {
		// Check if the drush executable exists in the current directory
		drushPath := filepath.Join(dir, drushExecutablePath)
		if _, err := os.Stat(drushPath); err == nil {
			return dir
		}

		// Move to the parent directory
		dir = filepath.Dir(dir)
	}

	// If no Drupal root is found, exit with an error message
	fmt.Println("Drupal root not found")
	os.Exit(1)
	return "" // This line is never reached, but it satisfies Go's return requirements
}

func determineDrushExecutable(drupalRoot string) string {
	// Build the path to the drush executable
	return filepath.Join(drupalRoot, drushExecutablePath)
}
