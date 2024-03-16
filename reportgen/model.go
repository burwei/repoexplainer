package reportgen

// Component represents a discovered component within the repository.
// This could be a struct, interface, function, etc., within a Go file.
type Component struct {
	File    string   `json:"file"`    // Path to the file where the component is defined
	Package string   `json:"package"` // Package name where the component is defined
	Name    string   `json:"name"`    // Name of the struct
	Type    string   `json:"type"`    // Component type (e.g., "struct", "interface", "function")
	Fields  []string `json:"fields"`  // Fields of the component (relevant for structs and interfaces)
	Methods []string `json:"methods"` // Methods attached to the component (relevant for structs and interfaces)
}

// ComponentMap maps a directory path to a slice of Components contained within.
// The key is constructed from the last two directory names (e.g., "/internal/repo").
type ComponentMap map[string]Component