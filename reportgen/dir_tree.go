package reportgen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileTraverser traverses the files in a directory tree starting from a root directory.
type FileTraverser struct {
	RootPath    string   // RootPath is the starting point for the traversal
	Files       []string // Files stores the paths of files found during traversal
	currentFile int      // currentFile tracks the current index in the Files slice
}

// NewFileTraverser creates a new FileTraverser for a given root directory.
func NewFileTraverser(rootPath string) *FileTraverser {
	ft := &FileTraverser{
		RootPath:    rootPath,
		currentFile: -1, // Start before the first element
	}
	ft.populateFiles()
	return ft
}

// populateFiles fills the Files slice with all file paths starting from RootPath.
func (ft *FileTraverser) populateFiles() {
	filepath.WalkDir(ft.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return fs.SkipDir // Skip hidden directories
			}
		} else {
			if strings.HasPrefix(d.Name(), ".") {
				return nil // Skip hidden files
			}

			ft.Files = append(ft.Files, path)
		}

		return nil
	})
}

// NextFile returns the next file in the traversal. When there are no more files, it returns false.
func (ft *FileTraverser) NextFile() (string, bool) {
	ft.currentFile++
	if ft.currentFile >= len(ft.Files) {
		return "", false // Indicates no more files are available
	}
	return ft.Files[ft.currentFile], true
}

// PrintDirectoryStructure prints the directory structure to the console.
func (ft *FileTraverser) PrintDirectoryStructure() (string, error) {
	if len(ft.Files) == 0 {
		return "", fmt.Errorf("no files have been traversed")
	}

	// Initialize directory structure map with root directory
	dirStructure := map[string][]string{
		".": {}, // Represent root directory with a dot, consistent with Unix-like filesystem notation
	}

	// Fill the directory structure map with files, organized by their directory paths
	for _, file := range ft.Files {
		relPath, err := filepath.Rel(ft.RootPath, file)
		if err != nil {
			return "", err
		}
		dir := filepath.Dir(relPath)
		base := filepath.Base(relPath)
		dirStructure[dir] = append(dirStructure[dir], base)
	}

	// Build the directory tree string
	var builder strings.Builder
	rootName := filepath.Base(ft.RootPath)
	builder.WriteString(fmt.Sprintf("/%s\n", rootName))
	// Generate sorted list of directories for consistent ordering
	dirs := make([]string, 0, len(dirStructure))
	for dir := range dirStructure {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		files := dirStructure[dir]
		sort.Strings(files) // Sort files for consistent order

		if dir != "." { // Skip the root directory since it's already added
			depth := strings.Count(dir, string(os.PathSeparator))
			indent := strings.Repeat("  ", depth)
			builder.WriteString(fmt.Sprintf("%s/%s\n", indent, filepath.Base(dir)))
		}

		for _, file := range files {
			depth := strings.Count(dir, string(os.PathSeparator)) + 1
			indent := strings.Repeat("  ", depth)
			builder.WriteString(fmt.Sprintf("%s- %s\n", indent, file))
		}
	}

	return builder.String(), nil
}
