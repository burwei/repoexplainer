package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/burwei/repoexplainer/compfinder"
	"github.com/burwei/repoexplainer/reportgen"
)

const (
	fileName = "repoexplain.md"
)

func Run(rootPath string) error {
	// Use the base name of the root directory as the repo name
	rootDirName := filepath.Base(rootPath)
	rg := reportgen.NewReportGenerator(rootDirName, rootPath, compfinder.NewFinderFactory())

	// Create a new file in the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %s", err)
	}

	file, err := os.Create(filepath.Join(cwd, fileName))
	if err != nil {
		return fmt.Errorf("creating report file: %s", err)
	}
	defer file.Close()

	err = rg.GenerateReport(file)
	if err != nil {
		return fmt.Errorf("generating report: %s", err)
	}

	return nil
}
