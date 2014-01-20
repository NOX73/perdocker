package perd

import (
	"bufio"
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
		secretEnd: "asdasdfasdfdas",
	}

	w.Start()
	return w
}

func (w *aworker) Start() {
	w.log("Starting", w.Lang.Name)

	go func() {
		w.clearContainer()
		w.startContainer()

		fileHost := w.tmpHost + w.Lang.ExecutableFile()
		fileGuest := w.tmpGuest + w.Lang.ExecutableFile()
		runCommand := w.Lang.RunCommand(fileGuest)

		for {
			c := <-w.in
			w.log("Precessing ...")

			ioutil.WriteFile(fileHost, []byte(c.Command()), 755)

			w.stdInOut.WriteString(runCommand + "\nEXITSTATUS=$?\necho " + w.secretEnd + "\necho $EXITSTATUS\necho " + w.secretEnd + " 1>&2\n")
			w.stdInOut.Flush()

			out := make([]byte, 0)
			er := make([]byte, 0)

			// Read stdOut
			for {
				line, _, _ := w.stdInOut.ReadLine()
				if string(line) == w.secretEnd {
					break
				}

				out = append(out, line...)
				out = append(out, eol)
			}

			// Read exitCode
			scode, _, _ := w.stdInOut.ReadLine()
			code, err := strconv.Atoi(string(scode))
			if err != nil {
				panic(err)
			}

			// Read stdErr
			for {
				line, _, _ := w.stdErr.ReadLine()
				if string(line) == w.secretEnd {
					break
				}

				er = append(er, line...)
				er = append(er, eol)
			}

			w.log("Finished ...")

			c.Response(out, er, code)
		}

		w.stopContainer()

	}()

}

func (w *aworker) startContainer() {
	container := exec.Command("docker", "run", "-i", "-v", w.tmpHost+":"+w.tmpGuest, "-name="+w.Name, w.Lang.Image, "/bin/bash", "-l")
	w.Container = container

	in, _ := container.StdinPipe()
	out, _ := container.StdoutPipe()
	er, _ := container.StderrPipe()

	w.stdInOut = bufio.NewReadWriter(bufio.NewReader(out), bufio.NewWriter(in))
	w.stdErr = bufio.NewReader(er)

	err := w.Container.Start()
	if err != nil {
		panic(err)
	}
}

func (w *aworker) stopContainer() {
	w.clearContainer()
}
