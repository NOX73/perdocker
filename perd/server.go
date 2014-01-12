package perd

import "net/http"

type Server interface {
  Run()
}

func NewServer () Server {
  ruby := NewRubyRunner()
  ruby.Run()

  return &server{ ruby }
}

type server struct {
  rubyRunner Runner
}

func (s *server) Run () {
  // Root path

  http.HandleFunc("/", s.rootHandler)

  http.ListenAndServe(":9000", nil)
}

func (s *server) rootHandler ( w http.ResponseWriter, r *http.Request ) {

  res := s.rubyRunner.Eval("puts 'Hello World'")
  w.Write(res.Bytes())

}


