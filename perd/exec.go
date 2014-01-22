package perd

import (
  "strconv"
  "time"
  "errors"
)

var (
  ErrExecTimout = errors.New("Execution timeout.")
)

type Exec struct {
  out <-chan []byte
  err <-chan []byte

  end []byte

  done chan bool

  StdOut []byte
  StdErr []byte
  ExitCode int
}

func NewExec (out, err <-chan []byte, end []byte) *Exec {
  return &Exec {
    out: out,
    err: err,
    end: end,
    done: make(chan bool, 1),
  }
}

func (e *Exec) Wait(timeout time.Duration) error {
  select {
  case <-e.done:
    return nil
  case <-time.After(timeout):
    return ErrExecTimout;
  }
}

func (e *Exec) Start () {

  for {
    if e.isFinish() {break}

    select {
    case line := <- e.out:
      if e.isEnd(line) {
        e.out = nil
        e.ExitCode = e.extractCode(line)
      } else {
        e.StdOut = append(e.StdOut, line...)
      }
    case line := <- e.err:
      if e.isEnd(line) {
        e.err = nil
      } else {
        e.StdErr = append(e.StdErr, line...)
      }
    }
  }

  e.done <- true
}

func (e *Exec) extractCode (line []byte) int {
  scode := string(line)[len(e.end) : len(line)-1]
  code, err := strconv.Atoi(scode)
  if err != nil { return 1 }
  return code
}

func (e *Exec) isEnd (line []byte) bool {
  return len(line) > len(e.end) && string(line[:len(e.end)]) == string(e.end)
}

func (e *Exec) isFinish () bool {
  return e.out == nil && e.err == nil
}
