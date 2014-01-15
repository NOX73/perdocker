package main

import "./perd"
import "sync"
import "flag"

var httpPort string

func main () {

  flag.StringVar(&httpPort, "port", "8080", "HTTP server port.")
  flag.Parse()

  var w sync.WaitGroup
  w.Add(1)

  server := perd.NewServer(httpPort)

  go func(){
    server.Run()
    w.Done()
  }()

  w.Wait()
}

