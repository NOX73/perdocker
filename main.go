package main

import (
	"./perd"
	"flag"
)

var httpPort string

func main() {
	var rubyWorkers int
	var nodejsWorkers int
	var golangWorkers int

	var timeout int64

	flag.StringVar(&httpPort, "port", "8080", "HTTP server port.")
	flag.IntVar(&rubyWorkers, "ruby-workers", 1, "Count of ruby workers.")
	flag.IntVar(&nodejsWorkers, "nodejs-workers", 1, "Count of nodejs workers.")
	flag.IntVar(&golangWorkers, "golang-workers", 1, "Count of go workers.")
	flag.Int64Var(&timeout, "timeout", 30, "Max execution time.")
	flag.Parse()

	workers := map[string]int{"ruby": rubyWorkers, "nodejs": nodejsWorkers, "golang": golangWorkers}
	server := perd.NewServer(httpPort, workers, timeout)

	server.Run()
}
