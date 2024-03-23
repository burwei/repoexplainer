package app

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/burwei/repoexplainer/compfinder"
	"github.com/burwei/repoexplainer/reportgen"
)

const (
	FileName = "repoexplain.md"
)

func Run(rootPath string, out io.Writer) error {
	// Use the base name of the root directory as the repo name
	rootDirName := filepath.Base(rootPath)
	rg := reportgen.NewReportGenerator(rootDirName, rootPath, compfinder.NewFinderFactory())

	err := rg.GenerateReport(out)
	if err != nil {
		return fmt.Errorf("generating report: %s", err)
	}

	return nil
}
