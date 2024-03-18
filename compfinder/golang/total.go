package golang

import (
	"strings"
	"sync"

	"github.com/burwei/repoexplainer/reportgen"
)

type ComponentFinder struct {
	structFinder       *StructFinder
	interfaceFinder    *InterfaceFinder
	funcFinder         *FuncFinder
	inMultiLineComment int
	inMultiLineString  bool
}

func NewComponentFinder() *ComponentFinder {
	return &ComponentFinder{
		structFinder:    NewStructFinder(),
		interfaceFinder: NewInterfaceFinder(),
		funcFinder:      NewFuncFinder(),
	}
}

func (cf *ComponentFinder) SetFile(filePath string) {
	cf.structFinder.SetFile(filePath)
	cf.interfaceFinder.SetFile(filePath)
	cf.funcFinder.SetFile(filePath)

	cf.inMultiLineComment = 0
	cf.inMultiLineString = false
}

func (cf *ComponentFinder) FindComponent(line string) {
	if strings.Contains(line, "/*") {
		cf.inMultiLineComment++
	}

	if strings.Contains(line, "*/") {
		cf.inMultiLineComment--
	}

	// The multiline string detection logic might not be perfect, but it's good enough most of the time.
	// We haven't considered (1) backticks inside single line comments, (2) escape characters
	if strings.Contains(line, "`") {
		var insideDoubleQuotes, insideSingleQuotes bool

		for _, char := range line {
			switch char {
			case '"':
				insideDoubleQuotes = !insideDoubleQuotes
			case '\'':
				insideSingleQuotes = !insideSingleQuotes
			case '`':
				if insideDoubleQuotes || insideSingleQuotes {
					// backtick inside quotes, so it's not the start or end of a multi-line string
					return // Exit early since we found a backtick inside quotes
				} else {
					// Found a backtick not inside quotes, toggle inMultiLineString and exit
					cf.inMultiLineString = !cf.inMultiLineString
					return
				}
			}
		}

		// backtick inside quotes, so it's not the start or end of a multi-line string
	}

	if cf.inMultiLineComment != 0 || cf.inMultiLineString {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		cf.structFinder.FindComponent(line)
		wg.Done()
	}()

	go func() {
		cf.interfaceFinder.FindComponent(line)
		wg.Done()
	}()

	go func() {
		cf.funcFinder.FindComponent(line)
		wg.Done()
	}()

	wg.Wait()
}

func (cf *ComponentFinder) GetComponents() reportgen.ComponentMap {
	components := reportgen.ComponentMap{}

	for key, val := range cf.structFinder.GetComponents() {
		components[key] = val
	}

	for key, val := range cf.interfaceFinder.GetComponents() {
		components[key] = val
	}

	for key, val := range cf.funcFinder.GetComponents() {
		structCompKey, dirPathBasedCompKey := cf.funcFinder.ConvertFuncCompKey(key)
		if structCompKey == "" {
			components[dirPathBasedCompKey] = val
		}

		if structComp, ok := components[structCompKey]; ok {
			// The function is a method of a struct, add it to the struct's methods
			components[structCompKey] = reportgen.Component{
				File:    structComp.File,
				Name:    structComp.Name,
				Package: structComp.Package,
				Type:    structComp.Type,
				Fields:  structComp.Fields,
				Methods: append(structComp.Methods, val.Name),
			}
		} else {
			// The function has a receiver, but the struct is not found
			// Usually this shouldn't happen, because this GetComponents() will be called
			// after all files are processed, and the struct should be found by then.
			// But, just in case, we add the function to the components map with the dirPathBasedCompKey
			components[dirPathBasedCompKey] = val
		}
	}

	return components
}
