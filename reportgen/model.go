package reportgen

// Component represents a discovered component within the repository.
// This could be a struct, interface, function, etc., within a Go file.
type Component struct {
	File    string   `json:"file"`    // Full path to the file where the component is defined
	Package string   `json:"package"` // Package name where the component is defined
	Name    string   `json:"name"`    // Name of the struct
	Type    string   `json:"type"`    // Component type (e.g., "struct", "interface", "function")
	Fields  []string `json:"fields"`  // Fields of the component (relevant for structs and interfaces)
	Methods []string `json:"methods"` // Methods attached to the component (relevant for structs and interfaces)
}

// ComponentMap maps a directory path to a slice of Components contained within.
// The key is a combination of the directory path and the component name.
// Key format: "path/to/dir:ComponentName". The path doesn't inclue the root directory and file name.
type ComponentMap map[string]Component

// OutputComponentMap maps a directory path to a slice of Components contained within.
// The key is the directory path.
// Key format: "path/to/dir". The path doesn't inclue the root directory and file name.
type OutputComponentMap map[string][]Component
