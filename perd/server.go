package perd

import "net/http"
import "io/ioutil"

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

  http.ListenAndServe(":1111", nil)
}

func (s *server) rootHandler ( w http.ResponseWriter, r *http.Request ) {

  body, err := ioutil.ReadAll(r.Body);

  if err == nil {
    res := s.rubyRunner.Eval(string(body))
    w.Write(res.Bytes())
  }

}


