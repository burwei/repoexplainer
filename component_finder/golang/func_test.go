package golang

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/burwei/repoexplainer/reportgen"
	"github.com/stretchr/testify/assert"
)

func TestFuncFinderFindComponent(t *testing.T) {
	testCases := []struct {
		name         string
		dirPath      string
		fileName     string
		fileContent  string
		expectedComp reportgen.ComponentMap
	}{
		{
			name:     "Simple function",
			dirPath:  "simple",
			fileName: "simple.go",
			fileContent: `
package simple

func SimpleFunc() int {
    return 1
}
`,
			expectedComp: reportgen.ComponentMap{
				":SimpleFunc": reportgen.Component{
					File:    filepath.Join("simple", "simple.go"),
					Package: "simple",
					Name:    "SimpleFunc() int",
					Type:    "func",
				},
			},
		},
		{
			name:     "Function with receiver",
			dirPath:  "complex",
			fileName: "complex.go",
			fileContent: `
package complex

type ComplexStruct struct {
    Value int
}

func (cs *ComplexStruct) GetValue() int {
    return cs.Value
}
`,
			expectedComp: reportgen.ComponentMap{
				"ComplexStruct:GetValue": reportgen.Component{
					File:    filepath.Join("complex", "complex.go"),
					Package: "complex",
					Name:    "GetValue() int",
					Type:    "func",
				},
			},
		},
		{
			name:     "Multiple functions in a file",
			dirPath:  "multi",
			fileName: "multi.go",
			fileContent: `
package multi

func FirstFunc() string {
    return "first"
}

func SecondFunc() int {
    return 2
}
`,
			expectedComp: reportgen.ComponentMap{
				":FirstFunc": reportgen.Component{
					File:    filepath.Join("multi", "multi.go"),
					Package: "multi",
					Name:    "FirstFunc() string",
					Type:    "func",
				},
				":SecondFunc": reportgen.Component{
					File:    filepath.Join("multi", "multi.go"),
					Package: "multi",
					Name:    "SecondFunc() int",
					Type:    "func",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ff := NewFuncFinder(tc.dirPath, tc.fileName, &reportgen.ComponentMap{})

			// Simulating line-by-line reading
			lines := strings.Split(tc.fileContent, "\n")
			for _, line := range lines {
				ff.FindComponent(line)
			}

			components := ff.GetFuncComponents()

			assert.Equal(t, tc.expectedComp, components)
		})
	}
}
