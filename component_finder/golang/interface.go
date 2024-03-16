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
	currentPkg       string
	currentInterface string
	fileName         string
	dirPath          string
}

func NewInterfaceFinder(dirPath, fileName string) *InterfaceFinder {
	return &InterfaceFinder{
		components: reportgen.ComponentMap{},
		dirPath:    dirPath,
		fileName:   fileName,
	}
}

func (ifd *InterfaceFinder) FindComponent(line string) {
	ifd.mu.Lock()
	defer ifd.mu.Unlock()

	if strings.HasPrefix(line, "package ") {
		ifd.currentPkg = strings.TrimSpace(line[len("package "):])
		return
	}

	// Interface definition or method detection logic
	if strings.Contains(line, "interface {") { // fast check
		if interfaceName := extractInterfaceName(line); interfaceName != "" { // detailed check
			compKey := getInterfaceCompKey(ifd.dirPath, interfaceName)
			ifd.currentInterface = interfaceName

			// In Go, there is only one interface with the same name in the same directory
			// So, we can ignore the duplicate interface definition
			if _, ok := ifd.components[compKey]; !ok {
				ifd.components[compKey] = reportgen.Component{
					File:    filepath.Join(ifd.dirPath, ifd.fileName),
					Package: ifd.currentPkg,
					Name:    interfaceName,
					Type:    "interface",
				}
			}
		}
	} else {
		if ifd.currentInterface != "" {
			// remove inline comments if any
			line = strings.Split(line, "//")[0]
			parts := strings.Fields(line)

			// close the interface definition if the line contains only "}"
			if len(parts) == 1 && parts[0] == "}" {
				ifd.currentInterface = ""

				return
			}

			// get the method definition
			method := strings.Join(parts, " ")
			compKey := getInterfaceCompKey(ifd.dirPath, ifd.currentInterface)

			ifd.components[compKey] = reportgen.Component{
				File:    filepath.Join(ifd.dirPath, ifd.fileName),
				Package: ifd.currentPkg,
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
	compCopy := make(reportgen.ComponentMap)
	for k, v := range ifd.components {
		compCopy[k] = v
	}

	return compCopy
}

func (ifd *InterfaceFinder) Close() {
	ifd.mu.Lock()
	defer ifd.mu.Unlock()

	ifd.currentPkg = ""
}

func getInterfaceCompKey(dirPath, interfaceName string) string {
	return dirPath + ":" + interfaceName
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
