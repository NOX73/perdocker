package perd

import "net/http"
import "io/ioutil"
import "log"

type Server interface {
  Run()
}

type config struct {
  port string
}

func NewServer (port string, workers map[string]int) Server {
  ruby := NewRubyRunner()
  ruby.RunWorkers(workers["ruby"])

  return &server{ &config{port} ,ruby }
}

type server struct {
  config *config
  rubyRunner Runner
}

func (s *server) Run () {
  // Root path

  http.HandleFunc("/ruby", s.rubyHandler)

  log.Println("Listen http on", s.config.port)
  http.ListenAndServe(":" + s.config.port, nil)
}

func (s *server) rubyHandler ( w http.ResponseWriter, r *http.Request ) {
  body, err := ioutil.ReadAll(r.Body);

  if err == nil {
    res := s.rubyRunner.Eval(string(body))
    w.Write(res.Bytes())
  }

}


