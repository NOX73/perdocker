package perd

import (
	"fmt"
)

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
	return fmt.Sprintf(l.Command, filePath)
}

// ExecutableFile returns filename, which will be used to store user's code
func (l *Lang) ExecutableFile() string {
	return l.FileName
}

// Ruby language settings
var Ruby = &Lang{"ruby", "run.rb", ".rb", "perdocker/ruby:attach", "ruby %s"}

// Nodejs settings
var Nodejs = &Lang{"nodejs", "index.js", ".js", "perdocker/nodejs:attach", "node %s"}

// Golang settings
var Golang = &Lang{"golang", "main.go", ".go", "perdocker/go:attach", "go run %s"}

// Python settings
var Python = &Lang{"python", "run.py", ".py", "perdocker/python:attach", "python %s"}

// C settings
var C = &Lang{"c", "a.c", ".c", "perdocker/c:attach", "gcc -o /tmp/a %s && /tmp/a"}

// CPP settings
var CPP = &Lang{"cpp", "a.cpp", ".cpp", "perdocker/c:attach", "g++ -o /tmp/a %s && /tmp/a"}

// PHP settings
var PHP = &Lang{"php", "index.php", ".php", "perdocker/php:attach", "php %s"}

// Universal languages container 
var Universal = &Lang{Image: "perdocker/universal:latest"}

var Languages = map[string]*Lang{
	"ruby":   Ruby,
	"nodejs": Nodejs,
	"golang": Golang,
	"python": Python,
	"c":      C,
	"cpp":    CPP,
	"php":    PHP,
}
