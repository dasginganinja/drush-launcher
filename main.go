package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	"github.com/dasginganinja/drush-launcher/drushlauncher"
)

var drupalRoot string

func main() {
	defaultRoot, _err := os.Getwd()

	if _err != nil {
		fmt.Println("Error getting current working directory:", _err)
		os.Exit(1)
	}

	// Strip program name from arguments before looping
	progArgs := os.Args[1:]

	for i, arg := range progArgs {
		if arg == "-r" || arg == "--root" {
			if i+1 < len(progArgs) {
				defaultRoot = progArgs[i+1]
			} else {
				fmt.Println("Error: Missing value for root argument")
				os.Exit(1)
			}
		} else if strings.HasPrefix(arg, "--root=") || strings.HasPrefix(arg, "-r=") {
			defaultRoot = strings.Split(arg,"=")[1]
		}
	}

	drupalRoot, _err = drushlauncher.FindDrupalRoot(defaultRoot)

	if _err != nil {
		fmt.Println(_err)
		os.Exit(1)
	}

	// Construct the full path to the drush executable
	drushExec := filepath.Join(drupalRoot, "vendor", "bin", "drush")

	// Check if the drush executable exists
	if _, err := os.Stat(drushExec); os.IsNotExist(err) {
		fmt.Println("Error: Drush executable not found at", drushExec)
		os.Exit(1)
	}

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
