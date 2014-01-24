package perd

import (
	"sync"
	"time"
)

const (
	killTimeout      = 5 * time.Second
	minWorkersCount  = 1
	newWorkerTimeout = 1 * time.Second
)

// Runner run several workers
type Runner interface {
	RunWorker()
	Eval(string) Result
}

type runner struct {
	Lang *Lang

	evalWorker chan Command
	newEval    chan Command
	killWorker chan bool

	workersCount    int64
	maxWorkersCount int64
	Timeout         int64
}

var workerID int64
var workerIDLock sync.Mutex

// NewRunner returns new Runner
func NewRunner(lang *Lang, workers int64, timeout int64) *runner {
	r := &runner{
		Lang: lang,

		evalWorker: make(chan Command),
		newEval:    make(chan Command),
		killWorker: make(chan bool, 1),

		maxWorkersCount: workers,
		Timeout:         timeout,
	}
	go r.Start()
}

func (r *runner) Start() {
	for {

		var killTimer <-chan time.Time
		if r.workersCount > minWorkersCount {
			killTimer = time.After(killTimeout)
		}

		select {
		case c := <-r.newEval:
			r.sendCommandToWorker(c)
		case <-killTimer:
			r.StopWorker()
		}

	}
}

func (r *runner) sendCommandToWorker(c Command) {

	var newWorkerTimer <-chan time.Time
	if r.workersCount < r.maxWorkersCount {
		newWorkerTimer = time.After(newWorkerTimeout)
	}

	select {
	case r.evalWorker <- c:
	case <-newWorkerTimer:
		r.RunWorker()
		r.evalWorker <- c
	}

}

func (r *runner) Eval(command string) Result {
	respCh := make(chan Result)
	r.newEval <- NewCommand(command, respCh)
	return <-respCh
}

func (r *runner) StopWorker() {

	select {
	case r.killWorker <- true:
		r.workersCount--
	default:
	}

}

func (r *runner) RunWorker() {
	workerIDLock.Lock()
	workerID++
	wid := workerID
	workerIDLock.Unlock()

	r.workersCount++

	NewWorker(r.Lang, wid, r.Timeout, r.evalWorker, r.killWorker)
}
