package reportgen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ReportGenerator struct {
	repoName      string
	fileTraverser *FileTraverser
	finderFactory FinderFactory
}

func NewReportGenerator(repoName, rootPath string, finderFactory FinderFactory) *ReportGenerator {
	return &ReportGenerator{
		repoName:      repoName,
		fileTraverser: NewFileTraverser(rootPath),
		finderFactory: finderFactory,
	}
}

func (rg *ReportGenerator) GenerateReport(out io.Writer) error {
	dirStructure, err := rg.fileTraverser.PrintDirectoryStructure()
	if err != nil {
		return fmt.Errorf("printing directory structure: %s", err)
	}

	err = rg.findCodeStructuresInFiles()
	if err != nil {
		return fmt.Errorf("finding code structures in files: %s", err)
	}

	outputCompMap := OutputComponentMap{}
	for _, finder := range rg.finderFactory.GetFinders() {
		compMap := finder.GetComponents()
		for key, comp := range compMap {
			dirPath := strings.Split(key, ":")[0]
			outputCompMap[dirPath] = append(outputCompMap[dirPath], comp)
		}
	}

	writer := bufio.NewWriter(out)
	writer.WriteString(fmt.Sprintf("# %s\n\n", rg.repoName))
	writer.WriteString("## directory structure\n\n")
	writer.WriteString(dirStructure)
	writer.WriteString("\n\n## components\n")

	for dirPath, comps := range outputCompMap {
		writer.WriteString(fmt.Sprintf(" - dir: %s\n", dirPath))
		for _, comp := range comps {
			writer.WriteString(fmt.Sprintf("     - %s\n", comp.Name))
			writer.WriteString(fmt.Sprintf("         - file: %s\n", comp.File))
			writer.WriteString(fmt.Sprintf("         - package: %s\n", comp.Package))
			writer.WriteString(fmt.Sprintf("         - type: %s\n", comp.Type))
			writer.WriteString("         - fields:\n")
			for _, field := range comp.Fields {
				writer.WriteString(fmt.Sprintf("             - %s\n", field))
			}
			writer.WriteString("         - methods:\n")
			for _, method := range comp.Methods {
				writer.WriteString(fmt.Sprintf("             - %s\n", method))
			}
		}
	}

	return nil
}

func (rg *ReportGenerator) findCodeStructuresInFiles() error {
	// iterate over all files in the repo
	filePath, ok := rg.fileTraverser.NextFile()
	for ok {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("opening file %s: %s", filePath, err)
		}
		defer file.Close()

		// Set the file for all the finders
		for _, finder := range rg.finderFactory.GetFinders() {
			finder.SetFile(filepath.Dir(filePath), filepath.Base(filePath))
		}

		// Loop through all the lines in the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			for _, finder := range rg.finderFactory.GetFinders() {
				finder.FindComponent(scanner.Text())
			}
		}

		// Check for errors during Scan. End of file is expected and not reported by Scan as an error.
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scanning file %s: %s", filePath, err)
		}

		filePath, ok = rg.fileTraverser.NextFile()
	}

	return nil
}
