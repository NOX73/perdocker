package perd

import (
  "sync"
)

type Runner interface {
  RunWorker()
  RunWorkers(int)
  Eval(string) Result
}

type runner struct {
  Lang      *Lang
  runCh     chan Command
  Timeout   int64
}

var workerId int64
var workerIdLock sync.Mutex

func NewRunner(lang *Lang, workers int, timeout int64) *runner {
  r := &runner{lang, make(chan Command), timeout}
  r.RunWorkers(workers)
  return r
}

func (r *runner) Eval (command string) Result {
  respCh := make(chan Result)
  r.runCh <- NewCommand(command, respCh)
  return <-respCh
}

func (r *runner) RunWorkers (count int) {
  for i := count; i>0; i-- {
    r.RunWorker()
  }
}

func (r *runner) RunWorker () {
  workerIdLock.Lock()
    workerId ++
    wid := workerId
  workerIdLock.Unlock()

  NewWorker(r.Lang, wid, r.Timeout, r.runCh )
}
