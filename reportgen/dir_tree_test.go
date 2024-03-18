package reportgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFiles(t *testing.T) {
	testCases := []struct {
		name          string
		directoryTree func(rootDir string) []string // Function to create a directory tree and return expected file paths
		expFileNames  []string                      // Expected file names (base names) found by FileTraverser
	}{
		{
			name: "Single level directory",
			directoryTree: func(rootDir string) []string {
				files := []string{"file1.txt", "file2.txt"}
				for _, file := range files {
					os.WriteFile(filepath.Join(rootDir, file), []byte("test"), 0644)
				}
				return files
			},
			expFileNames: []string{"file1.txt", "file2.txt"},
		},
		{
			name: "Nested directories with files",
			directoryTree: func(rootDir string) []string {
				nestedDir1 := filepath.Join(rootDir, "dir1")
				nestedDir2 := filepath.Join(nestedDir1, "dir2")
				os.MkdirAll(nestedDir2, 0755)

				files := []string{
					"file1.txt",
					filepath.Join("dir1", "file2.txt"),
					filepath.Join("dir1", "dir2", "file3.txt"),
				}
				expectedPaths := []string{}
				for _, file := range files {
					fullPath := filepath.Join(rootDir, file)
					os.WriteFile(fullPath, []byte("test"), 0644)
					expectedPaths = append(expectedPaths, filepath.Base(fullPath))
				}
				return expectedPaths
			},
			expFileNames: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			expectedFileNames := tc.directoryTree(tmpDir)

			ft := NewFileTraverser(tmpDir)

			// Assuming populateFiles is called inside NewFileTraverser
			foundFileNames := make([]string, len(ft.Files))
			for i, file := range ft.Files {
				foundFileNames[i] = filepath.Base(file)
			}

			assert.ElementsMatch(t, expectedFileNames, foundFileNames)
		})
	}
}

func TestNextFile(t *testing.T) {
	testCases := []struct {
		name          string
		directoryTree func(rootDir string) []string // Setup directories and files, return paths for verification
	}{
		{
			name: "Single level directory",
			directoryTree: func(rootDir string) []string {
				// Setup files
				files := []string{"a.txt", "b.txt"}
				for _, f := range files {
					path := filepath.Join(rootDir, f)
					os.WriteFile(path, []byte("test"), 0644)
				}
				return files
			},
		},
		{
			name: "Nested directories",
			directoryTree: func(rootDir string) []string {
				// Setup nested directory structure and files
				files := []string{
					"root.txt",
					filepath.Join("nested", "nested.txt"),
					filepath.Join("nested", "deeply", "deep.txt"),
				}
				for _, f := range files {
					fullPath := filepath.Join(rootDir, f)
					os.MkdirAll(filepath.Dir(fullPath), 0755)
					os.WriteFile(fullPath, []byte("test"), 0644)
				}
				return files
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tc.directoryTree(tmpDir)

			ft := NewFileTraverser(tmpDir)

			var foundFiles []string
			for {
				file, ok := ft.NextFile()
				if !ok {
					break
				}
				// We compare only the base names since the order should follow the slice population
				foundFiles = append(foundFiles, filepath.Base(file))
			}

			// Expected files are those returned by the directoryTree setup function,
			// adjusted to compare base names only, to match the foundFiles slice
			expectedFiles := []string{}
			for _, f := range tc.directoryTree(tmpDir) {
				expectedFiles = append(expectedFiles, filepath.Base(f))
			}

			assert.ElementsMatch(t, expectedFiles, foundFiles, "The found files do not match the expected files")
		})
	}
}

func TestPrintDirectoryStructure(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		setupDirectory func(rootDir string) // Function to set up directory and files
		expectedOutput string               // Expected string output of PrintDirectoryStructure
	}{
		{
			name: "Single level directory with files",
			setupDirectory: func(rootDir string) {
				// Create files in the root directory
				os.WriteFile(filepath.Join(rootDir, "file1.txt"), []byte("content"), 0644)
				os.WriteFile(filepath.Join(rootDir, "file2.txt"), []byte("content"), 0644)
			},
			expectedOutput: "/001\n\t- file1.txt\n\t- file2.txt\n",
		},
		{
			name: "Nested directories with files",
			setupDirectory: func(rootDir string) {
				// Create nested directory structure with files
				nestedDir := filepath.Join(rootDir, "dir1")
				os.Mkdir(nestedDir, 0755)
				os.WriteFile(filepath.Join(nestedDir, "nested_file1.txt"), []byte("content"), 0644)

				deeperNestedDir := filepath.Join(nestedDir, "dir2")
				os.Mkdir(deeperNestedDir, 0755)
				os.WriteFile(filepath.Join(deeperNestedDir, "nested_file2.txt"), []byte("content"), 0644)
			},
			expectedOutput: "/001\t\n/dir1\n\t\t- nested_file1.txt\n\t\t/dir2\n\t\t\t- nested_file2.txt\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tc.setupDirectory(tmpDir)

			ft := NewFileTraverser(tmpDir)
			output, err := ft.PrintDirectoryStructure()

			assert.NoError(t, err)
			// Normalize path separators for cross-platform compatibility
			expectedOutput := strings.ReplaceAll(tc.expectedOutput, "/", string(os.PathSeparator))
			assert.Equal(t, expectedOutput, output)
		})
	}
}
