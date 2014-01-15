package main

import (
  "./perd"
  "sync"
  "flag"
)

var httpPort string

func main () {
  var rubyWorkers int
  var w sync.WaitGroup
  w.Add(1)

  flag.StringVar(&httpPort, "port", "8080", "HTTP server port.")
  flag.IntVar(&rubyWorkers, "ruby-workers", 1, "HTTP server port.")
  flag.Parse()

  workers := map[string]int{ "ruby": rubyWorkers }
  server := perd.NewServer(httpPort, workers)

  go func(){
    server.Run()
    w.Done()
  }()

  w.Wait()
}

