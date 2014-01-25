package perd

// Lang is a struct to store language settings (name, Docker image, command to exec code, etc.).
type Lang struct {
	Name     string
	FileName string
	Ext      string
	Image    string
	Command  string
}

func (l *Lang) uniqFileName() string {
	return uniqFileName() + l.Ext
}

// RunCommand forms command using Lang.Command and a given string.
// Example: `ruby /tmp/perdocker/run.rb`
func (l *Lang) RunCommand(filePath string) string {
	return l.Command + " " + filePath
}

// ExecutableFile returns filename, which will be used to store user's code
func (l *Lang) ExecutableFile() string {
	return l.FileName
}

// Ruby language settings
var Ruby = &Lang{"ruby", "run.rb", ".rb", "perdocker/ruby:attach", "ruby"}

// Nodejs settings
var Nodejs = &Lang{"nodejs", "index.js", ".js", "perdocker/nodejs:attach", "node"}

// Golang settings
var Golang = &Lang{"golang", "main.go", ".go", "perdocker/go:attach", "go run"}

// Python settings
var Python = &Lang{"python", "run.py", ".py", "perdocker/python:attach", "python"}
