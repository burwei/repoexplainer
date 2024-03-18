package reportgen

// ComponentFinder is an interface for finding definitions of components like interfaces,
// structs, functions, etc., within files. It analyzes lines of code and identifies
// components based on the provided definitions.
type ComponentFinder interface {
	// SetFile sets the path of the current file being processed.
	// It's the beginning of a new file.
	SetFile(filePath string)

	// FindComponent takes a line of code as input and determines if it defines a component.
	// The input lines are continuous within a single file.
	// This method is responsible for parsing the code and extracting component definitions.
	FindComponent(line string)

	// GetComponents returns a ComponentMap of all components found by the finder.
	// This method allows retrieval of all identified components after processing
	// a file or a set of continuous lines of code, organized by their directory path.
	GetComponents() ComponentMap
}

// FinderFactory is an interface for creating ComponentFinder instances.
type FinderFactory interface {
	GetFinders() []ComponentFinder
}
