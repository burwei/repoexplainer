package compfinder

import (
	"github.com/burwei/repoexplainer/compfinder/golang"
	"github.com/burwei/repoexplainer/reportgen"
)

// FinderFactory is a struct that manages a collection of ComponentFinders.
type FinderFactory struct {
	Finders []reportgen.ComponentFinder
}

func NewFinderFactory() *FinderFactory {
	golangCompFinder := golang.NewComponentFinder()

	return &FinderFactory{
		Finders: []reportgen.ComponentFinder{golangCompFinder},
	}
}

// GetFinders creates instances of ComponentFinders and returns them.
func (ff *FinderFactory) GetFinders() []reportgen.ComponentFinder {
	return ff.Finders
}
