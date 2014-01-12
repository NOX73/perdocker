package main

import "./perd"
import "sync"

func main () {
  var w sync.WaitGroup
  w.Add(1)

  server := perd.NewServer()

  go func(){
    server.Run()
    w.Done()
  }()

  w.Wait()
}

