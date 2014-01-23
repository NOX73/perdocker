package perd

import (
	"log"
	"time"
)

const (
	eol byte = 10
)

type Worker interface {
	Start()
}

// This worker do not run container per request.
type worker struct {
	Lang       *Lang
	Id         int64
	MaxExecute time.Duration
	Container  Container

	in   chan Command
	exit chan bool
}

func NewWorker(lang *Lang, id, timeout int64, in chan Command, exit chan bool) (Worker, error) {

	container, err := NewContainer(id, lang)
	if err != nil {
		return nil, err
	}

	w := &worker{
		Container:  container,
		Lang:       lang,
		Id:         id,
		MaxExecute: time.Duration(timeout) * time.Second,

		in:   in,
		exit: exit,
	}

	go w.Start()
	return w, nil
}

func (w *worker) Start() {
	w.log("Starting ...")

	err := w.Container.Init()
	if err != nil {
		w.log("Can't start container.", err)
		return
	}

workerLoop:
	for {

		var c Command
		select {
		case c = <-w.in:
		case <-w.exit:
			break workerLoop
		}

		w.log("Precessing ...")

		var err error

		command := []byte(c.Command())
		exec, err := w.Container.Exec(command)

		err = exec.Wait(w.MaxExecute)

		if err != nil {
			w.log("Timeout kill. Restarting ...")
			c.Response(exec.StdOut, exec.StdErr, 137)

			// TODO: kill proccess instead restart container
			// it's required docker 0.8.0 feature for run command inside exists container.
			w.Container.Restart()

			continue
		}

		c.Response(exec.StdOut, exec.StdErr, exec.ExitCode)
		w.log("Finished ...")

		if w.Container.Clear() != nil {
			w.log("Container Clear error. Restarting ...")
			w.Container.Restart()
		}

	}

	w.Container.Stop()
	w.log("Stoping ...")

}

func (w *worker) log(s ...interface{}) {
	var params = make([]interface{}, 0)
	params = append(params, w.Lang.Name, "worker", w.Id, "\t")
	params = append(params, s...)
	log.Println(params...)
}
