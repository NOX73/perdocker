package perd

import (
  "strconv"
  "io/ioutil"
  "os/exec"
  "bytes"
  "syscall"
  "log"
  "time"
)

const (
  maxExecuteSeconds = 30
)

type Worker struct {
  Lang      *Lang
  Id        int64
  in        chan Command
  MaxExecute  time.Duration
}

func NewWorker (lang *Lang, id int64, in chan Command) *Worker {
  w := &Worker{lang, id, in, maxExecuteSeconds * time.Second}
  w.Start()
  return w
}

func (w *Worker) Start () {
  log.Println("Starting", w.Lang.Name, "worker ", w.Id)

  path := "/tmp/" + w.Lang.Name + "/"
  wName := "perdoker_" + w.Lang.Name +"_" + strconv.FormatInt(w.Id, 10)
  sharePath := path + ":" + path

  go func () {

    for {
      c := <- w.in
      log.Println("Worker", w.Id, ". Precessing", w.Lang.Name, "...")

      filePath := path + w.Lang.uniqFileName() 

      ioutil.WriteFile(filePath, []byte(c.Command()), 755)
      exec.Command("docker", "rm", wName).Run()

      cmd := exec.Command("docker", "run", "-v", sharePath, "-name=" + wName, w.Lang.Image, "/bin/bash", "-l", "-c", w.Lang.RunCommand(filePath))

      var stdOut, stdErr bytes.Buffer
      var code int
      cmd.Stdout, cmd.Stderr = &stdOut, &stdErr

      cmd.Start()

      done := make(chan error)
      go func () {
        done <- cmd.Wait()
      }()

      var err error
      select {
      case err = <- done:
      case <- time.After(w.MaxExecute):
        //cmd.Process.Kill()
        exec.Command("docker", "kill", wName).Run()
        log.Println("Worker", w.Id, ". Killed by timeout.")
        err = <- done
      }

      if err != nil { log.Println("Worker", w.Id, ". Error:", err) }

      code = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

      log.Println("Worker", w.Id, ". Code", code)

      c.Response(stdOut.Bytes(), stdErr.Bytes(), code)
    }

  }()

}
