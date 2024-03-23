package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/burwei/repoexplainer/app"
)

func main() {
	// Define a help flag
	helpFlag := flag.Bool("h", false, "Display help information")

	// Define a file output flag
	fileFlag := flag.Bool("f", false, "Write output to a file")

	flag.Parse()

	// Check if the help flag was provided
	if *helpFlag {
		fmt.Println("repoexplainer - Analyze a repository and generate a repoexlain.md at current directory")
		fmt.Println("\nUsage of repoexplainer:")
		fmt.Println("  repoexplainer [directory]")
		fmt.Println("  -h: Display help information")
		fmt.Println("  -f: Write output to a file")
		fmt.Println("\nExamples:")
		fmt.Println("  repoexplainer .                  # Analyze the current directory and copy output to clipboard")
		fmt.Println("  repoexplainer ./../another_repo  # Analyze a relative directory path and copy output to clipboard")
		fmt.Println("  repoexplainer /path/to/the/repo  # Analyze an absolute directory path and copy output to clipboard")
		fmt.Println("  repoexplainer -f .               # Analyze the current directory and write output to a file")
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

	// Write output to a file or copy to clipboard based on the flag
	if *fileFlag {
		// Write output to a file
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("getting current working directory: %s", err)
		}

		file, err := os.Create(filepath.Join(cwd, app.FileName))
		if err != nil {
			log.Fatalf("creating report file: %s", err)
		}
		defer file.Close()

		err = app.Run(absPath, file)
		if err != nil {
			log.Fatalf("Error running app: %s", err)
		}

		fmt.Println("Report file generated successfully!")
	} else {
		var buffer bytes.Buffer

		err := app.Run(absPath, &buffer)
		if err != nil {
			log.Fatalf("Error running app: %s", err)
		}

		err = clipboard.WriteAll(buffer.String())
		if err != nil {
			log.Fatalf("Error copying to clipboard: %s", err)
		}

		fmt.Println("Report copied to clipboard successfully!")
	}
}
