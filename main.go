package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dasginganinja/drush-launcher/drushlauncher"
)

func main() {
	// Parse command-line flags
	altRoot := ""

	// Strip program name from arguments before looping
	progArgs := os.Args[1:]

	for i, arg := range progArgs {
		// If we have -r or --root we will use the next argument as the drupal root (if exists)
		if arg == "-r" || arg == "--root" {
			if i+1 < len(progArgs) {
				altRoot = progArgs[i+1]
				break // no need to loop when we find the root
			} else {
				fmt.Println("Error: Missing value for root argument")
				os.Exit(1)
			}
		}
		// If we have --root={altRoot}
		if strings.HasPrefix(arg, "--root=") {
			altRoot = strings.Replace(arg, "--root=", "", 1)
			break // no need to loop when we find the root
		}
	}

	var drupalRoot string

	// Use the alternative Drupal root if provided
	if altRoot != "" {
		drupalRoot = altRoot
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
	// Parse the composer.json to get the bin-dir flag.
	// If no bin-dir flag is found, or no composer.json file use the default vendor/bin
	var drushRoot = filepath.Join("vendor", "bin")
	if (drushlauncher.GetComposerBinDir(drupalRoot) != "") {
		drushRoot = drushlauncher.GetComposerBinDir(drupalRoot)
	}

	drushExec := filepath.Join(drupalRoot, drushRoot, "drush")

	// Check if the drush executable exists
	if _, err := os.Stat(drushExec); os.IsNotExist(err) {
		fmt.Println("Error: Drush executable not found at", drushExec)
		os.Exit(1)
	}
	// Construct the full command to run drush
	drushCmd := exec.Command(drushExec, progArgs...)

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
