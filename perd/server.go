package perd

import "net/http"
import "io/ioutil"

type Server interface {
  Run()
}

type config struct {
  port string
}

func NewServer (port string) Server {
  ruby := NewRubyRunner()
  ruby.RunWorker()

  return &server{ &config{port} ,ruby }
}

type server struct {
  config *config
  rubyRunner Runner
}

func (s *server) Run () {
  // Root path

  http.HandleFunc("/ruby", s.rubyHandler)

  http.ListenAndServe(":"+s.config.port, nil)
}

func (s *server) rubyHandler ( w http.ResponseWriter, r *http.Request ) {
  body, err := ioutil.ReadAll(r.Body);

  if err == nil {
    res := s.rubyRunner.Eval(string(body))
    w.Write(res.Bytes())
  }

}


