package perd

import (
	"sync"
	"time"
)

const (
	killTimeout      = 5 * time.Second
	minWorkersCount  = 1
	newWorkerTimeout = 50 * time.Millisecond
	//Empirically chosen number in which there are no problems during normal operation
	systemDefaultProcessCount = 26
)

// Runner run several workers
type Runner interface {
	RunWorker()
	Eval(*Lang, string) Result
}

type runner struct {
	Lang *Lang

	evalWorker chan Command
	newEval    chan Command
	killWorker chan bool

	maxProcessCount int64

	workersCount    int64
	maxWorkersCount int64
	Timeout         int64
}

var workerID int64
var workerIDLock sync.Mutex

// NewRunner returns new Runner
func NewRunner(lang *Lang, workers int64, timeout int64) Runner {
	r := &runner{
		Lang: lang,

		evalWorker: make(chan Command),
		newEval:    make(chan Command),
		killWorker: make(chan bool, 1),

		maxProcessCount: workers * systemDefaultProcessCount,
		maxWorkersCount: workers,
		Timeout:         timeout,
	}
	go r.Start()
	return r
}

func (r *runner) Start() {

	for i := 0; i < minWorkersCount; i++ {
		r.RunWorker()
	}

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

func (r *runner) Eval(lang *Lang, code string) Result {
	respCh := make(chan Result)
	r.newEval <- NewCommand(lang, code, respCh)
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

	NewWorker(r.Lang, wid, r.Timeout, r.evalWorker, r.killWorker, r.maxProcessCount)
}
