package main

import (
  "./perd"
  "sync"
  "flag"
)

var httpPort string

func main () {
  var rubyWorkers, nodejsWorkers int
  var w sync.WaitGroup
  w.Add(1)

  flag.StringVar(&httpPort, "port", "8080", "HTTP server port.")
  flag.IntVar(&rubyWorkers, "ruby-workers", 1, "Count of ruby workers.")
  flag.IntVar(&nodejsWorkers, "nodejs-workers", 1, "Count of nodejs workers.")
  flag.Parse()

  workers := map[string]int{ "ruby": rubyWorkers, "nodejs": nodejsWorkers }
  server := perd.NewServer(httpPort, workers)

  go func(){
    server.Run()
    w.Done()
  }()

  w.Wait()
}

