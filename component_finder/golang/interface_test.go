package golang

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/burwei/repoexplainer/reportgen"
	"github.com/stretchr/testify/assert"
)

func TestInterfaceFinderFindComponent(t *testing.T) {
	testCases := []struct {
		name         string
		dirPath      string
		fileName     string
		fileContent  string
		expectedComp reportgen.ComponentMap
	}{
		{
			name:     "empty interface without methods",
			dirPath:  "empty",
			fileName: "empty.go",
			fileContent: `
package empty

type EmptyInterface interface {
}
`,
			expectedComp: reportgen.ComponentMap{
				"empty:EmptyInterface": reportgen.Component{
					File:    filepath.Join("empty", "empty.go"),
					Package: "empty",
					Name:    "EmptyInterface",
					Type:    "interface",
				},
			},
		},
		{
			name:     "Interface with multiple methods",
			dirPath:  "complex",
			fileName: "complex.go",
			fileContent: `
package complex

type ComplexInterface interface {
    GetName() string
    GetValue() int
}
`,
			expectedComp: reportgen.ComponentMap{
				"complex:ComplexInterface": reportgen.Component{
					File:    filepath.Join("complex", "complex.go"),
					Package: "complex",
					Name:    "ComplexInterface",
					Type:    "interface",
					Methods: []string{"GetName() string", "GetValue() int"},
				},
			},
		},
		{
			name:     "Multiple interfaces in a file",
			dirPath:  "multi",
			fileName: "multi.go",
			fileContent: `
package multi

type FirstInterface interface {
    GetFirstField() string
}

type SecondInterface interface {
    GetSecondField() int
}
`,
			expectedComp: reportgen.ComponentMap{
				"multi:FirstInterface": reportgen.Component{
					File:    filepath.Join("multi", "multi.go"),
					Package: "multi",
					Name:    "FirstInterface",
					Type:    "interface",
					Methods: []string{"GetFirstField() string"},
				},
				"multi:SecondInterface": reportgen.Component{
					File:    filepath.Join("multi", "multi.go"),
					Package: "multi",
					Name:    "SecondInterface",
					Type:    "interface",
					Methods: []string{"GetSecondField() int"},
				},
			},
		},
		{
			name:     "Data model interface with embedded interface",
			dirPath:  "models",
			fileName: "models.go",
			fileContent: `
package models

type BaseInterface interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type UserInterface interface {
	BaseInterface
	GetUsername() string
	GetEmail() string
}
		`,
			expectedComp: reportgen.ComponentMap{
				"models:BaseInterface": reportgen.Component{
					File:    filepath.Join("models", "models.go"),
					Package: "models",
					Name:    "BaseInterface",
					Type:    "interface",
					Methods: []string{"GetID() string", "GetCreatedAt() time.Time", "GetUpdatedAt() time.Time"},
				},
				"models:UserInterface": reportgen.Component{
					File:    filepath.Join("models", "models.go"),
					Package: "models",
					Name:    "UserInterface",
					Type:    "interface",
					Methods: []string{"BaseInterface", "GetUsername() string", "GetEmail() string"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ifd := NewInterfaceFinder(tc.dirPath, tc.fileName)

			// Simulating line-by-line reading
			lines := strings.Split(tc.fileContent, "\n")
			for _, line := range lines {
				ifd.FindComponent(line)
			}
			// Finalize parsing by calling Close
			ifd.Close()

			components := ifd.GetComponents()

			assert.Equal(t, tc.expectedComp, components)
		})
	}
}
