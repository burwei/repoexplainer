package compfinder

import (
	"github.com/burwei/repoexplainer/component_finder/golang"
	"github.com/burwei/repoexplainer/reportgen"
)

// FinderFactory is a struct that manages a collection of ComponentFinders.
type FinderFactory struct {
	Finders []reportgen.ComponentFinder
}

// GetFinders creates instances of ComponentFinders and returns them.
func (ff *FinderFactory) GetFinders(dirPath, fileName string) []reportgen.ComponentFinder {
	// golang
	structFinder := golang.NewStructFinder(dirPath, fileName)
	interfaceFinder := golang.NewInterfaceFinder(dirPath, fileName)
	funcFinder := golang.NewFuncFinder(dirPath, fileName)

	return []reportgen.ComponentFinder{structFinder, interfaceFinder, funcFinder}
}
