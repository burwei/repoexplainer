package golang

import (
	"sync"

	"github.com/burwei/repoexplainer/reportgen"
)

type ComponentFinder struct {
	structFinder    *StructFinder
	interfaceFinder *InterfaceFinder
	funcFinder      *FuncFinder
}

func NewComponentFinder() *ComponentFinder {
	return &ComponentFinder{
		structFinder:    NewStructFinder(),
		interfaceFinder: NewInterfaceFinder(),
		funcFinder:      NewFuncFinder(),
	}
}

func (cf *ComponentFinder) SetFile(dirPath, fileName string) {
	cf.structFinder.SetFile(dirPath, fileName)
	cf.interfaceFinder.SetFile(dirPath, fileName)
	cf.funcFinder.SetFile(dirPath, fileName)
}

func (cf *ComponentFinder) FindComponent(line string) {
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
