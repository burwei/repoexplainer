package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/burwei/repoexplainer/compfinder"
	"github.com/burwei/repoexplainer/reportgen"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %s", err)
	}

	// Use the base name of the current directory as the repo name
	repoName := filepath.Base(cwd)

	// Create a new report generator
	rg := reportgen.NewReportGenerator(repoName, cwd, compfinder.NewFinderFactory())

	// Create a new file in the current directory
	file, err := os.Create(filepath.Join(cwd, "repoexplainer.md"))
	if err != nil {
		log.Fatalf("Error creating report file: %s", err)
	}
	defer file.Close()

	// Generate the report
	err = rg.GenerateReport(file)
	if err != nil {
		log.Fatalf("Error generating report: %s", err)
	}
}
