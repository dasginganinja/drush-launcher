package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"flag"

	"github.com/dasginganinja/drush-launcher/drushlauncher"
)

var drupalRoot string

func main() {
	defaultRoot, _err := os.Getwd()
	const usage = "path to the Drupal root directory. Defaults to current directory."

	if _err != nil {
		fmt.Println("Error getting current working directory:", _err)
		os.Exit(1)
	}

	flag.StringVar(&drupalRoot, "root", defaultRoot, usage)
	flag.StringVar(&drupalRoot, "r", defaultRoot, usage + " (shorthand)")

	flag.Parse()

	drupalRoot, _err = drushlauncher.FindDrupalRoot(drupalRoot)

	if _err != nil {
		fmt.Println(_err)
		os.Exit(1)
	}


	// Construct the full path to the drush executable
	// Parse the composer.json to get the bin-dir flag.
	// If no bin-dir flag is found, use the default vendor/bin
	drushExec := filepath.Join(drupalRoot, drushlauncher.GetComposerBinDir(drupalRoot), "drush")

	// Check if the drush executable exists
	if _, err := os.Stat(drushExec); os.IsNotExist(err) {
		fmt.Println("Error: Drush executable not found at", drushExec)
		os.Exit(1)
	}
	
	// Construct the full command to run drush
	fmt.Println("Running drush with arguments:", flag.Args())
	flag.Set("root", drupalRoot)

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
