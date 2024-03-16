package golang

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/burwei/repoexplainer/reportgen"
	"github.com/stretchr/testify/assert"
)

func TestStructFinderFindComponent(t *testing.T) {
	testCases := []struct {
		name         string
		dirPath      string
		fileName     string
		fileContent  string
		expectedComp reportgen.ComponentMap
	}{
		{
			name:     "Simple struct without methods",
			dirPath:  "simple",
			fileName: "simple.go",
			fileContent: `
package simple

type SimpleStruct struct {
    ID int
}
`,
			expectedComp: reportgen.ComponentMap{
				"simple:SimpleStruct": reportgen.Component{
					File:    filepath.Join("simple", "simple.go"),
					Package: "simple",
					Name:    "SimpleStruct",
					Type:    "struct",
					Fields:  []string{"ID int"},
				},
			},
		},
		{
			name:     "Struct with multiple fields and methods",
			dirPath:  "complex",
			fileName: "complex.go",
			fileContent: `
package complex

type ComplexStruct struct {
    Name  string
    Value int
}

func (cs *ComplexStruct) GetName() string {
    return cs.Name
}

func (cs *ComplexStruct) GetValue() int {
    return cs.Value
}
`,
			expectedComp: reportgen.ComponentMap{
				"complex:ComplexStruct": reportgen.Component{
					File:    filepath.Join("complex", "complex.go"),
					Package: "complex",
					Name:    "ComplexStruct",
					Type:    "struct",
					Fields:  []string{"Name string", "Value int"},
				},
			},
		},
		{
			name:     "Multiple structs in a file",
			dirPath:  "multi",
			fileName: "multi.go",
			fileContent: `
package multi

type FirstStruct struct {
    FirstField string
}

type SecondStruct struct {
    SecondField int
}

func (fs *FirstStruct) GetFirstField() string {
    return fs.FirstField
}

func (ss *SecondStruct) GetSecondField() int {
    return ss.SecondField
}
`,
			expectedComp: reportgen.ComponentMap{
				"multi:FirstStruct": reportgen.Component{
					File:    filepath.Join("multi", "multi.go"),
					Package: "multi",
					Name:    "FirstStruct",
					Type:    "struct",
					Fields:  []string{"FirstField string"},
				},
				"multi:SecondStruct": reportgen.Component{
					File:    filepath.Join("multi", "multi.go"),
					Package: "multi",
					Name:    "SecondStruct",
					Type:    "struct",
					Fields:  []string{"SecondField int"},
				},
			},
		},
		{
			name:     "Data model struct with embedded field",
			dirPath:  "models",
			fileName: "models.go",
			fileContent: `
package models

type BaseModel struct {
    ID        string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type User struct {
    BaseModel
    Username string
    Email    string
}
`,
			expectedComp: reportgen.ComponentMap{
				"models:BaseModel": reportgen.Component{

					File:    filepath.Join("models", "models.go"),
					Package: "models",
					Name:    "BaseModel",
					Type:    "struct",
					Fields:  []string{"ID string", "CreatedAt time.Time", "UpdatedAt time.Time"},
				},
				"models:User": reportgen.Component{
					File:    filepath.Join("models", "models.go"),
					Package: "models",
					Name:    "User",
					Type:    "struct",
					Fields:  []string{"BaseModel", "Username string", "Email string"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sf := NewStructFinder(tc.dirPath, tc.fileName)

			// Simulating line-by-line reading
			lines := strings.Split(tc.fileContent, "\n")
			for _, line := range lines {
				sf.FindComponent(line)
			}
			// Finalize parsing by calling Close
			sf.FileEnd()

			components := sf.GetComponents()

			assert.Equal(t, tc.expectedComp, components)
		})
	}
}
