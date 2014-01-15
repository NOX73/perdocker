package perd

import "net/http"
import "io/ioutil"

type Server interface {
  Run()
}

func NewServer () Server {
  ruby := NewRubyRunner()
  ruby.RunWorker()

  return &server{ ruby }
}

type server struct {
  rubyRunner Runner
}

func (s *server) Run () {
  // Root path

  http.HandleFunc("/ruby", s.rubyHandler)

  http.ListenAndServe(":1111", nil)
}

func (s *server) rubyHandler ( w http.ResponseWriter, r *http.Request ) {
  body, err := ioutil.ReadAll(r.Body);

  if err == nil {
    res := s.rubyRunner.Eval(string(body))
    w.Write(res.Bytes())
  }

}


