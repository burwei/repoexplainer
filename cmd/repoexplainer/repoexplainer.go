package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/burwei/repoexplainer/app"
)

func main() {
	// Define a help flag
	helpFlag := flag.Bool("h", false, "Display help information")
	flag.Parse()

	// Check if the help flag was provided
	if *helpFlag {
		fmt.Println("repoexplainer - Analyze a repository and generate a repoexlain.md at current directory")
		fmt.Println("\nUsage of repoexplainer:")
		fmt.Println("  repoexplainer [directory]")
		fmt.Println("  -h: Display help information")
		fmt.Println("\nExamples:")
		fmt.Println("  repoexplainer .                  # Analyze the current directory")
		fmt.Println("  repoexplainer ./../another_repo  # Analyze a relative directory path")
		fmt.Println("  repoexplainer /path/to/the/repo  # Analyze an absolute directory path")
		return
	}

	var dirPath string

	// Check if the user has provided a directory path as an argument
	if len(flag.Args()) > 0 {
		// Use the provided directory path
		dirPath = flag.Arg(0)
	} else {
		// Default to the current working directory if no argument is provided
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %s", err)
		}
		dirPath = cwd
	}

	if !filepath.IsAbs(dirPath) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not get current working directory: %s", err)
		}
		dirPath = filepath.Join(cwd, dirPath)
	}

	// Clean up the path to resolve any ".." or "." segments
	absPath := filepath.Clean(dirPath)

	// Create a new file in the current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getting current working directory: %s", err)
	}

	file, err := os.Create(filepath.Join(cwd, app.FileName))
	if err != nil {
		log.Fatalf("creating report file: %s", err)
	}
	defer file.Close()

	// Run the application with the absolute directory path
	err = app.Run(absPath, file)
	if err != nil {
		log.Fatalf("Error running app: %s", err)
	}

	fmt.Println("Report generated successfully!")
}
