package perd

import (
  "strconv"
  "io/ioutil"
  "os/exec"
  "bytes"
  "syscall"
  "log"
  "time"
  "os"
)

const (
  maxExecuteSeconds = 60
)

type Worker struct {
  Lang      *Lang
  Id        int64
  in        chan Command
  MaxExecute  time.Duration
  Name      string
  tmpHost   string
  tmpGuest  string
}

func NewWorker (lang *Lang, id, timeout int64, in chan Command) *Worker {

  wName := "perdoker_" + lang.Name +"_" + strconv.FormatInt(id, 10)
  tmpHostPath := "/tmp/perdocker/" + lang.Name + "/" + wName + "/"
  tmpGuestPath := "/tmp/perdocker/"

  err := os.MkdirAll(tmpHostPath, 0755)
  if err != nil { log.Println(err) }

  if timeout > maxExecuteSeconds { timeout = maxExecuteSeconds }

  w := &Worker{lang, id, in, time.Duration(timeout) * time.Second, wName, tmpHostPath, tmpGuestPath }
  w.Start()
  return w
}

func (w *Worker) Start () {
  w.log("Starting", w.Lang.Name)

  go func () {

    w.clearContainer()

    fileHost := w.tmpHost + w.Lang.ExecutableFile()
    fileGuest := w.tmpGuest + w.Lang.ExecutableFile()
    runCommand := w.Lang.RunCommand(fileGuest)

    for {
      c := <- w.in
      w.log( "Precessing", w.Lang.Name, "...")

      ioutil.WriteFile(fileHost, []byte(c.Command()), 755)

      cmd := exec.Command("docker", "run", "-v", w.tmpHost + ":" + w.tmpGuest, "-name=" + w.Name, w.Lang.Image, "/bin/bash", "-l", "-c", runCommand)

      var stdOut, stdErr bytes.Buffer
      var code int

      cmd.Stdout, cmd.Stderr = &stdOut, &stdErr

      cmd.Start()

      done := make(chan error)
      go func () {
        done <- cmd.Wait()
      }()

      select {
      case <- done:
        code = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
      case <- time.After(w.MaxExecute):
        w.log("Killed by timeout")
        w.clearContainer()

        // manualy set killed status
        code = 137

        <- done
      }

      w.log("Exit status code: ", code)

      c.Response(stdOut.Bytes(), stdErr.Bytes(), code)

      w.clearContainer()
    }

  }()

}

func (w *Worker) log (s ...interface{}) {
  var params = make([]interface{}, 0)
  params = append(params, "Worker:", w.Id, ">")
  params = append(params, s...)
  log.Println(params...)
}

func (w *Worker) killContainer () error {
  return exec.Command("docker", "kill", w.Name).Run()
}

func (w *Worker) rmContainer () error {
  return exec.Command("docker", "rm", w.Name).Run()
}

func (w *Worker) clearContainer () {
  for w.containerExist() {
    w.killContainer()
    w.rmContainer()
  }
}

func (w *Worker) containerExist () bool {
  err := exec.Command("docker", "inspect", w.Name).Run()

  return err == nil
}

