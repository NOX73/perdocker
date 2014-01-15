package main

import "./perd"
import "sync"
import "flag"

func main () {

  var port = flag.String("port", "8080", "HTTP server port.")


  var w sync.WaitGroup
  w.Add(1)

  server := perd.NewServer(*port)

  go func(){
    server.Run()
    w.Done()
  }()

  w.Wait()
}

