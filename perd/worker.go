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
  maxExecuteSeconds = 60
)

type Worker struct {
  Lang      *Lang
  Id        int64
  in        chan Command
  MaxExecute  time.Duration
  Name      string
  Path      string
  SharePath string
}

func NewWorker (lang *Lang, id, timeout int64, in chan Command) *Worker {

  wName := "perdoker_" + lang.Name +"_" + strconv.FormatInt(id, 10)
  path := "/tmp/" + lang.Name + "/"
  sharePath := path + ":" + path

  if timeout > maxExecuteSeconds { timeout = maxExecuteSeconds }

  w := &Worker{lang, id, in, time.Duration(timeout) * time.Second, wName, path, sharePath }
  w.Start()
  return w
}

func (w *Worker) Start () {
  w.log("Starting", w.Lang.Name)

  go func () {

    w.clearContainer()

    for {
      c := <- w.in
      w.log( "Precessing", w.Lang.Name, "...")

      filePath := w.Path + w.Lang.uniqFileName() 

      ioutil.WriteFile(filePath, []byte(c.Command()), 755)
      // eval code
      cmd := exec.Command("docker", "run", "-v", w.SharePath, "-name=" + w.Name, w.Lang.Image, "/bin/bash", "-l", "-c", w.Lang.RunCommand(filePath))

      var stdOut, stdErr bytes.Buffer
      var code int
      cmd.Stdout, cmd.Stderr = &stdOut, &stdErr

      cmd.Start()

      done := make(chan error)
      go func () {
        done <- cmd.Wait()
      }()

      // wait timout
      var err error
      select {
      case err = <- done:
      case <- time.After(w.MaxExecute):
        // TODO: not so cool
        //cmd.Process.Kill()
        w.killContainer()

        w.log("Killed by timeout")
        err = <- done
      }

      if err == nil { 
        code = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
      } else {
        code = err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus()
      }

      w.log("Code", code)

      c.Response(stdOut.Bytes(), stdErr.Bytes(), code)

      w.rmContainer()
    }

  }()

}

func (w *Worker) log (s ...interface{}) {
  var params = make([]interface{}, 0)
  params = append(params, "Worker:", w.Id, ">")
  params = append(params, s...)
  log.Println(params...)
}

func (w *Worker) killContainer () {
  // kill while container run
  for { if exec.Command("docker", "kill", w.Name).Run() != nil {break} else {w.log("Can't kill")} }
}

func (w *Worker) rmContainer () {
  // remove while container exist
  for { if exec.Command("docker", "rm", w.Name).Run() != nil {break} else {w.log("Can't rm")} }
}

func (w *Worker) clearContainer () {
  w.killContainer()
  w.rmContainer()
}
