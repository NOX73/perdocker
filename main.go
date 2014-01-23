package main

import (
	"./perd"
	"flag"
)

var httpListen = flag.String("listen", ":8080", "HTTP server bind to address & port. Ex: localhost:80 or :80")

var rubyWorkers = flag.Int64("ruby-workers", 1, "Count of ruby workers.")
var nodejsWorkers = flag.Int64("nodejs-workers", 1, "Count of nodejs workers.")
var golangWorkers = flag.Int64("golang-workers", 1, "Count of golang workers.")

var timeout = flag.Int64("timeout", 30, "Max execution time.")

func main() {
	flag.Parse()

	workers := map[string]int64{"ruby": *rubyWorkers, "nodejs": *nodejsWorkers, "golang": *golangWorkers}
	server := perd.NewServer(*httpListen, workers, *timeout)

	server.Run()
}
