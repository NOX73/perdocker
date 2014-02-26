package perd

import (
	"errors"
	"strconv"
	"time"
)

var (
	// ErrExecTimout occurs when the command runs out of maximum execution time
	ErrExecTimout = errors.New("execution timeout")

	// ErrReadStds throwed if std read has been terminated
	ErrReadStds = errors.New("std read terminated")
)

// Exec is a struct which collects stdout/stderr outputs from commands.
type Exec struct {
	out <-chan []byte
	err <-chan []byte

	end []byte

	done chan error

	StdOut   []byte
	StdErr   []byte
	ExitCode int
}

// NewExec returns new Exec
func NewExec(out, err <-chan []byte, end []byte) *Exec {
	e := &Exec{
		out:  out,
		err:  err,
		end:  end,
		done: make(chan error, 1),
	}
	go e.Start()
	return e
}

// Wait returns ErrExecTimout error if command runs out maximum execution time.
func (e *Exec) Wait(timeout time.Duration) error {
	select {
	case err := <-e.done:
		return err
	case <-time.After(timeout):
		return ErrExecTimout
	}
}

// Start collects stdout/stderr output and extract exit code.
func (e *Exec) Start() {

	for {
		if e.isFinish() {
			break
		}

		select {
		case line, ok := <-e.out:
			if !ok {
				e.done <- ErrReadStds
				return
			}

			if e.isEnd(line) {
				e.out = nil
				e.ExitCode = e.extractCode(line)
				e.StdOut = cutLast(e.StdOut)
			} else {
				e.StdOut = append(e.StdOut, line...)
			}

		case line, ok := <-e.err:
			if !ok {
				e.done <- ErrReadStds
				return
			}

			if e.isEnd(line) {
				e.err = nil
				e.StdErr = cutLast(e.StdErr)
			} else {
				e.StdErr = append(e.StdErr, line...)
			}

		}
	}

	e.done <- nil
}

func (e *Exec) extractCode(line []byte) int {
	scode := string(line)[len(e.end) : len(line)-1]
	code, err := strconv.Atoi(scode)
	if err != nil {
		return 1
	}
	return code
}

func (e *Exec) isEnd(line []byte) bool {
	return len(line) > len(e.end) && string(line[:len(e.end)]) == string(e.end)
}

func (e *Exec) isFinish() bool {
	return e.out == nil && e.err == nil
}

func cutLast(sl []byte) []byte {
	size := len(sl)
	if size > 0 {
		return sl[:size-1]
	}
	return sl
}
