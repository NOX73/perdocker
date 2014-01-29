package main

import (
	"./perd"
	"flag"
)

var httpListen = flag.String("listen", ":8080", "HTTP server bind to address & port. Ex: localhost:80 or :80")

var rubyWorkers = flag.Int64("ruby-workers", 1, "Count of ruby workers.")
var nodejsWorkers = flag.Int64("nodejs-workers", 1, "Count of nodejs workers.")
var golangWorkers = flag.Int64("golang-workers", 1, "Count of golang workers.")
var pythonWorkers = flag.Int64("python-workers", 1, "Count of python workers.")
var cWorkers = flag.Int64("c-workers", 1, "Count of C workers.")
var cppWorkers = flag.Int64("cpp-workers", 1, "Count of C++ workers.")
var phpWorkers = flag.Int64("php-workers", 1, "Count of PHP workers.")

var timeout = flag.Int64("timeout", 30, "Max execution time.")

var separate = flag.Bool("separate", false, "Separate workers by languages.")
var workers = flag.Int64("workers", 1, "Count of workers for non separated workers.")

func main() {
	flag.Parse()

	var server perd.Server

	if *separate {
		server = separatedServer()
	} else {
		server = nonSeparatedServer()
	}

	server.Run()
}

func nonSeparatedServer() perd.Server {
	return perd.NewUniversalServer(*httpListen, *workers, *timeout)
}

func separatedServer() perd.Server {

	workersMap := map[string]int64{
		"ruby":   *rubyWorkers,
		"nodejs": *nodejsWorkers,
		"golang": *golangWorkers,
		"python": *pythonWorkers,
		"c":      *cWorkers,
		"cpp":    *cppWorkers,
		"php":    *phpWorkers,
	}

	return perd.NewServer(*httpListen, workersMap, *timeout)
}
