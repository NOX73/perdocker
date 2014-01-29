package perd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// Server is a simple http server who listens for incoming requests.
// When request comes, he evals its body using Runner
type Server interface {
	Run()
}

type config struct {
	listen string
}

// NewServer returns new server
func NewServer(listen string, workers map[string]int64, timeout int64) Server {
	jsRunner := NewRunner(Nodejs, workers["nodejs"], timeout)
	runners := map[string]Runner{
		"ruby":       NewRunner(Ruby, workers["ruby"], timeout),
		"nodejs":     jsRunner,
		"javascript": jsRunner,
		"golang":     NewRunner(Golang, workers["golang"], timeout),
		"python":     NewRunner(Python, workers["python"], timeout),
		"c":          NewRunner(C, workers["c"], timeout),
		"cpp":        NewRunner(CPP, workers["cpp"], timeout),
		"php":        NewRunner(PHP, workers["php"], timeout),
	}
	return &server{&config{listen}, runners}
}

// NewServer returns new server with universal runner
func NewUniversalServer(listen string, workers, timeout int64) Server {
	runner := NewRunner(Universal, workers, timeout)
	runners := map[string]Runner{
		"ruby":       runner,
		"nodejs":     runner,
		"javascript": runner,
		"golang":     runner,
		"python":     runner,
		"c":          runner,
		"cpp":        runner,
		"php":        runner,
	}
	return &server{&config{listen}, runners}
}

type server struct {
	config  *config
	runners map[string]Runner
}

var (
	ErrUndefinedLang  = errors.New("Undefined Language.")
	ErrCantFindRunner = errors.New("Can't find runner")
)

func (s *server) Run() {
	// Root path

	http.HandleFunc("/api/evaluate", s.evaluateHandler)

	http.HandleFunc("/api/evaluate/ruby", s.rubyHandler)
	http.HandleFunc("/api/evaluate/nodejs", s.nodejsHandler)
	http.HandleFunc("/api/evaluate/javascript", s.nodejsHandler)
	http.HandleFunc("/api/evaluate/golang", s.golangHandler)
	http.HandleFunc("/api/evaluate/python", s.pythonHandler)
	http.HandleFunc("/api/evaluate/c", s.cHandler)
	http.HandleFunc("/api/evaluate/cpp", s.cppHandler)
	http.HandleFunc("/api/evaluate/php", s.phpHandler)

	log.Println("Listen http on", s.config.listen)
	http.ListenAndServe(s.config.listen, nil)
}

func (s *server) langHandler(w http.ResponseWriter, r *http.Request, lang *Lang) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	res, err := s.eval(lang, string(body))

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(res.Bytes())
}

func (s *server) nodejsHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, Nodejs)
}

func (s *server) rubyHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, Ruby)
}

func (s *server) golangHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, Golang)
}

func (s *server) pythonHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, Python)
}

func (s *server) cHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, C)
}

func (s *server) cppHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, CPP)
}

func (s *server) phpHandler(w http.ResponseWriter, r *http.Request) {
	s.langHandler(w, r, PHP)
}

type RequestJson struct {
	Lang string `json:"language"`
	Code string `json:"code"`
}

func (s *server) evaluateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var body []byte
	var res Result

	body, err = ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return
	}

	js := &RequestJson{}
	err = json.Unmarshal(body, js)

	if err != nil {
		log.Println(err)
		return
	}

	lang, ok := Languages[js.Lang]

	if !ok {
		log.Println(ErrUndefinedLang)
		return
	}

	res, err = s.eval(lang, js.Code)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(res.Bytes())
}

func (s *server) eval(lang *Lang, code string) (Result, error) {
	runner, ok := s.runners[lang.Name]
	if !ok {
		return nil, ErrCantFindRunner
	}
	return runner.Eval(lang, code), nil
}
