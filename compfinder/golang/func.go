package golang

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/burwei/repoexplainer/reportgen"
)

// FuncFinder is a ComponentFinder implementation for finding function definitions within Go files.
type FuncFinder struct {
	mu          sync.Mutex
	components  reportgen.ComponentMap
	fileName    string
	dirPath     string
	packageName string
}

func NewFuncFinder() *FuncFinder {
	return &FuncFinder{
		components: reportgen.ComponentMap{},
	}
}

// SetFile sets the directory path and file name for the current file being processed.
// It's the beginning of a new file.
func (ff *FuncFinder) SetFile(dirPath, fileName string) {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	ff.dirPath = dirPath
	ff.fileName = fileName
	ff.packageName = ""
}

func (ff *FuncFinder) FindComponent(line string) {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	if strings.HasPrefix(line, "package ") {
		ff.packageName = strings.TrimSpace(line[len("package "):])
		return
	}

	// Function definition detection logic
	if strings.HasPrefix(line, "func ") {
		funcSignature, receiver := extractFuncSignature(line)
		if funcSignature != "" {
			compKey := getFuncCompKey(receiver, funcSignature)

			// In Go, there can't be multiple functions with the same name with same receiver type
			// So, we don't need to handle duplicate function definitions
			ff.components[compKey] = reportgen.Component{
				File:    ff.fileName,
				Package: ff.packageName,
				Name:    funcSignature,
				Type:    "func",
			}
		}
	}
}

func (ff *FuncFinder) GetComponents() reportgen.ComponentMap {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	// Return a copy of the map to avoid race conditions
	// when the caller iterates over the map
	compCopy := make(reportgen.ComponentMap)
	for k, v := range ff.components {
		compCopy[k] = v
	}

	return compCopy
}

func (ff *FuncFinder) ConvertFuncCompKey(compKey string) (string, string) {
	parts := strings.Split(compKey, ":")
	comp := ff.components[compKey]
	structCompKey := filepath.Dir(comp.File) + ":" + parts[0]
	funcName := strings.Split(comp.Name, "(")[0]
	dirPathBasedCompKey := filepath.Dir(comp.File) + ":" + funcName

	// receiver part of the compKey is empty
	if parts[0] == "" {
		return "", dirPathBasedCompKey
	}

	return structCompKey, dirPathBasedCompKey
}

func getFuncCompKey(receiver, funcSignature string) string {
	funcName := strings.Split(funcSignature, "(")[0]
	if receiver == "" {
		return ":" + funcName
	}

	return receiver + ":" + funcName
}

// Assumes method signatures line follows the pattern "func (r ReceiverType) MethodName() ReturnType {".
// or "func MethodName() ReturnType {".
func extractFuncSignature(line string) (string, string) {
	parts := strings.Fields(line)

	if parts[0] != "func" {
		// it's not a method definition
		return "", ""
	}

	// function with receiver
	// follows the pattern "func (r ReceiverType) MethodName() ReturnType {"
	if parts[1][0] == '(' {
		// extract the string between the first ") " and " {"
		methodSignature := strings.SplitN(line, ") ", 2)[1]
		methodSignature = strings.Split(methodSignature, " {")[0]

		receiverStructType := strings.Split(parts[2], ")")[0]
		receiverStructType = strings.ReplaceAll(receiverStructType, "*", "")

		return methodSignature, receiverStructType
	}

	// function without receiver
	// follows the pattern "func MethodName() ReturnType {"
	// extract the string between the "func " and " {"
	methodSignature := strings.Split(line, "func ")[1]
	methodSignature = strings.Split(methodSignature, " {")[0]

	return methodSignature, ""
}
