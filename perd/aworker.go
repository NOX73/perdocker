package perd

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	eol byte = 10
)

// This worker do not run container per request.
type aworker struct {
	*worker

	secretEnd string

	Container *exec.Cmd
	stdInOut  *bufio.ReadWriter
	stdErr    *bufio.Reader

	outChan chan []byte
	errChan chan []byte

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewAWorker(lang *Lang, id, timeout int64, in chan Command) Worker {

	wName := "perdoker_" + lang.Name + "_" + strconv.FormatInt(id, 10)
	tmpHostPath := "/tmp/perdocker/" + lang.Name + "/" + wName + "/"
	tmpGuestPath := "/tmp/perdocker/"

	err := os.MkdirAll(tmpHostPath, 0755)
	if err != nil {
		log.Println(err)
	}

	if timeout > maxExecuteSeconds {
		timeout = maxExecuteSeconds
	}

	w := &aworker{
		worker:    &worker{lang, id, in, time.Duration(timeout) * time.Second, wName, tmpHostPath, tmpGuestPath},
		secretEnd: "raabbbccc",
		outChan:   make(chan []byte, 5),
		errChan:   make(chan []byte, 5),
	}

	w.Start()
	return w
}

func (w *aworker) Start() {
	w.log("Starting ...")

	go func() {
		w.clearContainer()
		w.startContainer()

		fileHost := w.tmpHost + w.Lang.ExecutableFile()
		fileGuest := w.tmpGuest + w.Lang.ExecutableFile()
		runCommand := w.Lang.RunCommand(fileGuest)

		for {
			c := <-w.in
			w.log("Precessing ...")

			ioutil.WriteFile(fileHost, []byte(c.Command()), 0755)

			w.stdInOut.WriteString(runCommand + " 3>&- \n")
			w.stdInOut.WriteString("echo " + w.secretEnd + "$?\n")
			w.stdInOut.WriteString("echo " + w.secretEnd + " 1>&2\n")

			w.stdInOut.Flush()

			out := make([]byte, 0)
			er := make([]byte, 0)
			var code int
			var err error

			outChan := w.outChan
			errChan := w.errChan

			timeout := time.After(w.MaxExecute)

			for {

				if errChan == nil && outChan == nil {
					break
				}

				select {
				case line := <-outChan:

					if string(line)[:len(w.secretEnd)] == w.secretEnd {
						scode := string(line)[len(w.secretEnd) : len(line)-1]
						code, err = strconv.Atoi(scode)

						if err != nil {
							code = 1
						}

						outChan = nil
					} else {
						out = append(out, line...)
					}

				case line := <-errChan:

					if string(line)[:len(w.secretEnd)] == w.secretEnd {
						errChan = nil
					} else {
						er = append(er, line...)
					}

				case <-timeout:
					w.log("Timeout kill ...")
					errChan = nil
					outChan = nil
					code = 137
				}
			}

			w.log("Finished ...")
			c.Response(out, er, code)

      // TODO: kill proccess instead restart container
      // it's required docker 0.8.0 feature for run command inside exists container.
			if code != 137 {
				w.checkContainer()
			} else {
				w.restartContainer()
			}

		}

		w.stopContainer()
		w.log("Stoping ...")

	}()

}

func (w *aworker) checkContainer() {
  //TODO: Fork detector
}

func (w *aworker) startContainer() {
	container := exec.Command("docker", "run", "-m", "10m", "-c", "1", "-i", "-v", w.tmpHost+":"+w.tmpGuest+":ro", "-name="+w.Name, w.Lang.Image, "/bin/bash", "-l")
	w.Container = container

	w.stdin, _ = container.StdinPipe()
	w.stdout, _ = container.StdoutPipe()
	w.stderr, _ = container.StderrPipe()

	w.stdInOut = bufio.NewReadWriter(bufio.NewReader(w.stdout), bufio.NewWriter(w.stdin))
	w.stdErr = bufio.NewReader(w.stderr)

	err := w.Container.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			line, err := w.stdInOut.ReadBytes(eol)
			if err != nil {
				break
			}
			w.outChan <- line
		}
		w.log("StdInOut closed.", err)
	}()

	go func() {
		for {
			line, err := w.stdErr.ReadBytes(eol)
			if err != nil {
				break
			}
			w.errChan <- line
		}
		w.log("StdErr read closed.", err)
	}()
}

func (w *aworker) stopContainer() {
	w.stdin.Close()
	w.stdout.Close()
	w.stderr.Close()

	w.clearContainer()
}

func (w *aworker) restartContainer() {
	w.stopContainer()
	w.startContainer()
}
