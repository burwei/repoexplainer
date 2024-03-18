package golang

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/burwei/repoexplainer/reportgen"
)

// InterfaceFinder is a ComponentFinder implementation for finding interface definitions within Go files.
type InterfaceFinder struct {
	mu               sync.Mutex
	components       reportgen.ComponentMap
	currentInterface string
	filePath         string
	packageName      string
}

func NewInterfaceFinder() *InterfaceFinder {
	return &InterfaceFinder{
		components: reportgen.ComponentMap{},
	}
}

// SetFile sets the directory path and file name for the current file being processed.
// It's the beginning of a new file.
func (ifd *InterfaceFinder) SetFile(filePath string) {
	ifd.mu.Lock()
	defer ifd.mu.Unlock()

	ifd.filePath = filePath
	ifd.packageName = ""
	ifd.currentInterface = ""
}

func (ifd *InterfaceFinder) FindComponent(line string) {
	ifd.mu.Lock()
	defer ifd.mu.Unlock()

	if strings.HasPrefix(line, "package ") {
		ifd.packageName = strings.TrimSpace(line[len("package "):])
		return
	}

	// Interface definition or method detection logic
	if strings.Contains(line, "interface {") { // fast check
		if interfaceName := extractInterfaceName(line); interfaceName != "" { // detailed check
			compKey := getInterfaceCompKey(ifd.filePath, interfaceName)
			ifd.currentInterface = interfaceName

			// In Go, there is only one interface with the same name in the same directory
			// So, we can ignore the duplicate interface definition
			if _, ok := ifd.components[compKey]; !ok {
				ifd.components[compKey] = reportgen.Component{
					File:    ifd.filePath,
					Package: ifd.packageName,
					Name:    interfaceName,
					Type:    "interface",
				}
			}
		}
	} else {
		if ifd.currentInterface != "" {
			parts := strings.Fields(line)

			// close the interface definition if the line contains only "}"
			if len(parts) == 1 && parts[0] == "}" {
				ifd.currentInterface = ""

				return
			}

			// get the method definition
			method := strings.Join(parts, " ")
			if method == "" {
				return
			}

			compKey := getInterfaceCompKey(ifd.filePath, ifd.currentInterface)

			ifd.components[compKey] = reportgen.Component{
				File:    ifd.filePath,
				Package: ifd.packageName,
				Name:    ifd.currentInterface,
				Type:    "interface",
				Methods: append(ifd.components[compKey].Methods, method),
			}
		}
	}
}

func (ifd *InterfaceFinder) GetComponents() reportgen.ComponentMap {
	ifd.mu.Lock()
	defer ifd.mu.Unlock()

	// Return a copy of the map to avoid race conditions
	// when the caller iterates over the map
	compCopy := make(reportgen.ComponentMap)
	for k, v := range ifd.components {
		compCopy[k] = v
	}

	return compCopy
}

func getInterfaceCompKey(filePath, interfaceName string) string {
	return filepath.Dir(filePath) + ":" + interfaceName
}

// Assumes the interface declaration line follows the pattern "type InterfaceName interface {".
func extractInterfaceName(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		// it's not an interface definition
		return ""
	}
	if parts[0] == "type" && parts[2] == "interface" {
		return parts[1]
	}

	return ""
}
