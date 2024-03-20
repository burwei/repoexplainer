package reportgen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File struct {
	Type string
	Path string
}

// FileTraverser traverses the files in a directory tree starting from a root directory.
type FileTraverser struct {
	RootPath    string              // RootPath is the starting point for the traversal
	Files       []File              // Files stores the paths of files found during traversal
	existsFiles map[string]struct{} // existsFiles is a set of file paths to check for existence
	currentFile int                 // currentFile tracks the current index in the Files slice
}

// NewFileTraverser creates a new FileTraverser for a given root directory.
func NewFileTraverser(rootPath string) *FileTraverser {
	ft := &FileTraverser{
		RootPath:    rootPath,
		currentFile: -1, // Start before the first element
		existsFiles: make(map[string]struct{}),
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

			if _, ok := ft.existsFiles[path]; !ok {
				ft.Files = append(ft.Files, File{Type: "dir", Path: path})
				ft.existsFiles[path] = struct{}{}
			}
		} else {
			if strings.HasPrefix(d.Name(), ".") {
				return nil // Skip hidden files
			}

			ft.Files = append(ft.Files, File{Type: "file", Path: path})
			ft.existsFiles[path] = struct{}{}
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

	// Skip directories
	for ft.Files[ft.currentFile].Type == "dir" {
		ft.currentFile++
		if ft.currentFile >= len(ft.Files) {
			return "", false // Indicates no more files are available
		}
	}

	return ft.Files[ft.currentFile].Path, true
}

// PrintDirectoryStructure prints the directory structure to the console.
func (ft *FileTraverser) PrintDirectoryStructure() (string, error) {
	if len(ft.Files) == 0 {
		return "", fmt.Errorf("no files have been traversed")
	}

	// Initialize directory structure map with root directory
	dirStructure := map[string][]File{}

	// Fill the directory structure map with files, organized by their directory paths
	for _, file := range ft.Files {
		dir := filepath.Dir(file.Path)
		dirStructure[dir] = append(dirStructure[dir], File{Type: file.Type, Path: file.Path})
	}

	var builder strings.Builder
	dirs := make([]string, 0, len(dirStructure))
	for dir := range dirStructure {
		dirs = append(dirs, dir)
	}

	// Generate sorted list of directories for consistent ordering
	sort.Strings(dirs)

	offset := strings.Count(ft.RootPath, string(os.PathSeparator))
	for _, dir := range dirs {
		depth := strings.Count(dir, string(os.PathSeparator)) - offset
		if depth == -1 {
			continue // Skip the parent of root directory
		}

		indent := strings.Repeat("\t", depth)
		builder.WriteString(fmt.Sprintf("%s/%s\n", indent, filepath.Base(dir)))

		files := dirStructure[dir]
		for _, file := range files {

			if file.Type == "dir" {
				// Also print the empty directories
				// These directories are not present in the file list
				// so they won't be the key in the dirStructure map
				if _, ok := dirStructure[file.Path]; !ok {
					depth := strings.Count(dir, string(os.PathSeparator)) - offset + 1
					indent := strings.Repeat("\t", depth)
					builder.WriteString(fmt.Sprintf("%s/%s\n", indent, filepath.Base(file.Path)))

					continue
				}

				continue
			}

			depth := strings.Count(dir, string(os.PathSeparator)) - offset + 1
			indent := strings.Repeat("\t", depth)
			builder.WriteString(fmt.Sprintf("%s- %s\n", indent, filepath.Base(file.Path)))
		}
	}

	return builder.String(), nil
}
