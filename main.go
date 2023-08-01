package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/dasginganinja/drush-launcher/drushlauncher"
)

func main() {
	// Parse command-line flags
	altRoot := flag.String("r", "", "Set an alternative Drupal root")
	altRootLong := flag.String("root", "", "Set an alternative Drupal root (long form)")
	flag.Parse()

	var drupalRoot string
	var err error

	// Use the alternative Drupal root if provided
	if *altRoot != "" {
		drupalRoot = *altRoot
	} else if *altRootLong != "" {
		drupalRoot = *altRootLong
	} else {
		// If no alternative root provided, find the Drupal root from the current directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
		drupalRoot, err = drushlauncher.FindDrupalRoot(cwd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Construct the full path to the drush executable
	drushExec := filepath.Join(drupalRoot, "vendor", "bin", "drush")

	// Check if the drush executable exists
	if _, err := os.Stat(drushExec); os.IsNotExist(err) {
		fmt.Println("Error: Drush executable not found at", drushExec)
		os.Exit(1)
	}

	// Construct the full command to run drush
	drushCmd := exec.Command(drushExec, flag.Args()...)

	// Pass the current environment variables to the drush command
	drushCmd.Env = os.Environ()

	// Set the correct working directory for the drush command
	drushCmd.Dir = drupalRoot

	// Redirect standard input/output/error for the drush command
	drushCmd.Stdin = os.Stdin
	drushCmd.Stdout = os.Stdout
	drushCmd.Stderr = os.Stderr

	// Run the drush command
	if err := drushCmd.Run(); err != nil {
		fmt.Println("Error executing drush:", err)
		os.Exit(1)
	}
}
