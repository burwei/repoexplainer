package golang

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/burwei/repoexplainer/reportgen"
)

// FuncFinder is a ComponentFinder implementation for finding function definitions within Go files.
type FuncFinder struct {
	mu         sync.Mutex
	components reportgen.ComponentMap
	currentPkg string
	fileName   string
	dirPath    string
}

func NewFuncFinder(dirPath, fileName string, compMap *reportgen.ComponentMap) *FuncFinder {
	if compMap == nil {
		compMap = &reportgen.ComponentMap{}
	}

	return &FuncFinder{
		components: *compMap,
		dirPath:    dirPath,
		fileName:   fileName,
	}
}

func (ff *FuncFinder) FindComponent(line string) {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	if strings.HasPrefix(line, "package ") {
		ff.currentPkg = strings.TrimSpace(line[len("package "):])
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
				File:    filepath.Join(ff.dirPath, ff.fileName),
				Package: ff.currentPkg,
				Name:    funcSignature,
				Type:    "func",
			}
		}
	}
}

func (ff *FuncFinder) GetFuncComponents() reportgen.ComponentMap {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	// Return a copy of the map to avoid race conditions
	compCopy := make(reportgen.ComponentMap)
	for k, v := range ff.components {
		compCopy[k] = v
	}

	return compCopy
}

func (ff *FuncFinder) Close() {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	ff.currentPkg = ""
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
