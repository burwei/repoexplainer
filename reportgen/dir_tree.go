package reportgen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	TypeFile = "file"
	TypeDir  = "dir"
)

type File struct {
	Type string
	Path string
}

// FileTraverser traverses the files in a directory tree starting from a root directory.
type FileTraverser struct {
	RootPath    string // RootPath is the starting point for the traversal
	Files       []File // Files stores the paths of files found during traversal
	currentFile int    // currentFile tracks the current index in the Files slice
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
		if strings.HasPrefix(d.Name(), ".") { // Skip hidden files and directories
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		fileType := TypeFile
		if d.IsDir() {
			fileType = TypeDir
		}
		ft.Files = append(ft.Files, File{Type: fileType, Path: path})
		return nil
	})
}

// NextFile returns the next file in the traversal. When there are no more files, it returns false.
func (ft *FileTraverser) NextFile() (string, bool) {
	ft.currentFile++

	// Skip directories
	for ft.currentFile < len(ft.Files) && ft.Files[ft.currentFile].Type == TypeDir {
		ft.currentFile++
	}

	if ft.currentFile < len(ft.Files) {
		return ft.Files[ft.currentFile].Path, true
	}

	return "", false
}

// PrintDirectoryStructure prints the directory structure to the console.
func (ft *FileTraverser) PrintDirectoryStructure() (string, error) {
	if len(ft.Files) == 0 {
		return "", fmt.Errorf("no files have been traversed")
	}

	// Create a map of directories to files to maintain the structure
	dirStructure := map[string][]File{}
	for _, file := range ft.Files {
		dir := filepath.Dir(file.Path)
		dirStructure[dir] = append(dirStructure[dir], File{Type: file.Type, Path: file.Path})
	}

	dirs := make([]string, 0, len(dirStructure))
	for dir := range dirStructure {
		dirs = append(dirs, dir)
	}

	// Generate sorted list of directories for consistent ordering
	sort.Strings(dirs)

	var builder strings.Builder
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

			if file.Type == TypeDir {
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
