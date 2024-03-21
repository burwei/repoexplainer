package golang

import (
	"strings"
	"testing"

	"github.com/burwei/repoexplainer/reportgen"
	"github.com/stretchr/testify/assert"
)

func TestStructFinderFindComponent(t *testing.T) {
	testCases := []struct {
		name         string
		filePath     string
		fileContent  string
		expectedComp reportgen.ComponentMap
	}{
		{
			name:     "Simple struct without methods",
			filePath: "simple/simple.go",
			fileContent: `
package simple

type SimpleStruct struct {
    ID int
}
`,
			expectedComp: reportgen.ComponentMap{
				"simple:SimpleStruct": reportgen.Component{
					File:    "simple/simple.go",
					Package: "simple",
					Name:    "SimpleStruct",
					Type:    TypeStruct,
					Fields:  []string{"ID int"},
				},
			},
		},
		{
			name:     "Struct with multiple fields and methods",
			filePath: "complex/complex.go",
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
					File:    "complex/complex.go",
					Package: "complex",
					Name:    "ComplexStruct",
					Type:    TypeStruct,
					Fields:  []string{"Name string", "Value int"},
				},
			},
		},
		{
			name:     "Multiple structs in a file",
			filePath: "multi/multi.go",
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
					File:    "multi/multi.go",
					Package: "multi",
					Name:    "FirstStruct",
					Type:    TypeStruct,
					Fields:  []string{"FirstField string"},
				},
				"multi:SecondStruct": reportgen.Component{
					File:    "multi/multi.go",
					Package: "multi",
					Name:    "SecondStruct",
					Type:    TypeStruct,
					Fields:  []string{"SecondField int"},
				},
			},
		},
		{
			name:     "Data model struct with embedded field",
			filePath: "models/models.go",
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

					File:    "models/models.go",
					Package: "models",
					Name:    "BaseModel",
					Type:    TypeStruct,
					Fields:  []string{"ID string", "CreatedAt time.Time", "UpdatedAt time.Time"},
				},
				"models:User": reportgen.Component{
					File:    "models/models.go",
					Package: "models",
					Name:    "User",
					Type:    TypeStruct,
					Fields:  []string{"BaseModel", "Username string", "Email string"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sf := NewStructFinder()
			sf.SetFile(tc.filePath)

			// Simulating line-by-line reading
			lines := strings.Split(tc.fileContent, "\n")
			for _, line := range lines {
				sf.FindComponent(line)
			}

			components := sf.GetComponents()

			assert.Equal(t, tc.expectedComp, components)
		})
	}
}
