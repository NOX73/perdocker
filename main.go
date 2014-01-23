package main

import (
	"./perd"
	"flag"
)

func main() {

	var httpListen = *flag.String("listen", ":8080", "HTTP server bind to address & port. Ex: localhost:80 or :80")

	var rubyWorkers = *flag.Int("ruby-workers", 1, "Count of ruby workers.")
	var nodejsWorkers = *flag.Int("nodejs-workers", 1, "Count of nodejs workers.")
	var golangWorkers = *flag.Int("golang-workers", 1, "Count of golang workers.")

	var timeout = *flag.Int("timeout", 30, "Max execution time.")

	flag.Parse()

	workers := map[string]int{"ruby": rubyWorkers, "nodejs": nodejsWorkers, "golang": golangWorkers}
	server := perd.NewServer(httpListen, workers, int64(timeout))

	server.Run()
}
