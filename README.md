# repoexplainer
Generate a Markdown file to describe an existing repo, so that developers can explain their whole repo to an AI chatbot more easily.  

Currently, it only supports Go.   
However, other languages could be added easily.    

As a heavy ChatGPT-4 user, I use it a lot when I'm programming.  
Quite often, I need to describe my local/private repo to it for it to understand what I'm doing.   
I hope this tool will make it easier for developers to communicate with chat-based AI.   

## Installation
```
go install github.com/burwei/repoexplainer/cmd/repoexplainer@latest
```

## How to generate repoexplain.md
To analyze the repo at current directory 
```
repoexplainer
```
To analyze a repo at some other directory 
```
// relative path
repoexplainer ./../another_repo

// absolute path
repoexplainer /path/to/some/other/repo
```
Then, the report file "repoexplain.md" will be generated in the current directory.  

## How to use repoexplain.md
Here are some useful prompts I frequently use:  
```
Hi, I'm working on a repo and I need you help.
Here's the overview of the repo:

<upload/paste the repoexplain.md here>
```
```
Impelement XXX for me in YYY file.
```
```
Write me a xxx_mock.go file using testify/mock.
```
```
Write me the xxx_test.go file in table driven testing pattern with at least 3 test cases in each test function.
```
```
Add XXX function in YYY file to ZZZ.
```

## What does the markdown file look like
It looks like this:
```
# repoexplainer

## directory structure

/repoexplainer
	- README.md
	- go.mod
	- go.sum
	/app
		- app.go
	/cmd
		/repoexplainer
			- repoexplainer.go
	/compfinder
		- finder_factory.go
		/golang
			- const.go
			- func.go
			- func_test.go
			- interface.go
			- interface_test.go
			- struct.go
			- struct_test.go
			- total.go
			- total_test.go
	/example
		- repoexplain.md
	/reportgen
		- dir_tree.go
		- dir_tree_test.go
		- generator.go
		- interface.go
		- model.go


## components
 - dir: /repoexplainer/reportgen
     - ReportGenerator
         - file: /Users/burwei/Documents/github/repoexplainer/repoexplainer/reportgen/generator.go
         - package: reportgen
         - type: struct
         - fields:
             - rootDirName string
             - rootPath string
             - fileTraverser *FileTraverser
             - finderFactory FinderFactory
         - methods:
             - getOutputCompMap() OutputComponentMap
             - findCodeStructuresInFiles() error
             - GenerateReport(out io.Writer) error
...
```

To check the complete repoexplain.md example: [link](https://github.com/burwei/repoexplainer/blob/main/example/repoexplain.md)
