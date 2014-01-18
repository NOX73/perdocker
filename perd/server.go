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

func NewServer(port string, workers map[string]int, timeout int64) Server {
	ruby := NewRunner(Ruby, workers["ruby"], timeout)
	nodejs := NewRunner(Nodejs, workers["nodejs"], timeout)

	return &server{&config{port}, ruby, nodejs}
}

type server struct {
	config       *config
	rubyRunner   Runner
	nodejsRunner Runner
}

func (s *server) Run() {
	// Root path

	http.HandleFunc("/ruby", s.rubyHandler)
	http.HandleFunc("/nodejs", s.nodejsHandler)

	log.Println("Listen http on", s.config.port)
	http.ListenAndServe(":"+s.config.port, nil)
}

func (s *server) nodejsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		res := s.nodejsRunner.Eval(string(body))
		w.Write(res.Bytes())
	}

}

func (s *server) rubyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		res := s.rubyRunner.Eval(string(body))
		w.Write(res.Bytes())
	}

}
