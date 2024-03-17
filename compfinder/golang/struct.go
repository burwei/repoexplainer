package golang

import (
	"strings"
	"sync"

	"github.com/burwei/repoexplainer/reportgen"
)

// StructFinder is a ComponentFinder implementation for finding struct definitions within Go files.
// Not including methods.
type StructFinder struct {
	mu            sync.Mutex
	components    reportgen.ComponentMap
	currentStruct string
	fileName      string
	dirPath       string
	packageName   string
}

func NewStructFinder() *StructFinder {
	return &StructFinder{
		components: reportgen.ComponentMap{},
	}
}

// SetFile sets the directory path and file name for the current file being processed.
// It's the beginning of a new file.
func (sf *StructFinder) SetFile(dirPath, fileName string) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	sf.dirPath = dirPath
	sf.fileName = fileName
	sf.packageName = ""
	sf.currentStruct = ""
}

func (sf *StructFinder) FindComponent(line string) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	if strings.HasPrefix(line, "package ") {
		sf.packageName = strings.TrimSpace(line[len("package "):])
		return
	}

	// Struct definition or field detection logic
	if strings.Contains(line, "struct {") { // fast check
		if structName := extractStructName(line); structName != "" { // detailed check
			compKey := getStructCompKey(sf.dirPath, structName)
			sf.currentStruct = structName

			// In Go, there is only one struct with the same name in the same directory
			// So, we can ignore the duplicate struct definition
			if _, ok := sf.components[compKey]; !ok {
				sf.components[compKey] = reportgen.Component{
					File:    sf.fileName,
					Package: sf.packageName,
					Name:    structName,
					Type:    "struct",
				}
			}
		}
	} else {
		if sf.currentStruct != "" {
			// remove inline comments if any
			line = strings.Split(line, "//")[0]
			parts := strings.Fields(line)

			// close the struct definition if the line contains only "}"
			if len(parts) == 1 && parts[0] == "}" {
				sf.currentStruct = ""

				return
			}

			// get the field definition
			field := strings.Join(parts, " ")
			if field == "" {
				return
			}

			compKey := getStructCompKey(sf.dirPath, sf.currentStruct)

			sf.components[compKey] = reportgen.Component{
				File:    sf.fileName,
				Package: sf.packageName,
				Name:    sf.currentStruct,
				Type:    "struct",
				Fields:  append(sf.components[compKey].Fields, field),
				Methods: sf.components[compKey].Methods,
			}
		}
	}
}

func (sf *StructFinder) GetComponents() reportgen.ComponentMap {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	// Return a copy of the map to avoid race conditions
	// when the caller iterates over the map
	compCopy := make(reportgen.ComponentMap)
	for k, v := range sf.components {
		compCopy[k] = v
	}

	return compCopy
}

func getStructCompKey(dirPath, structName string) string {
	return dirPath + ":" + structName
}

// Assumes the struct declaration line follows the pattern "type StructName struct {".
func extractStructName(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		// it's not a struct definition
		return ""
	}
	if parts[0] == "type" && parts[2] == "struct" {
		return parts[1]
	}

	return ""
}
