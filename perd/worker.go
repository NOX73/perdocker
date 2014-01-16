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
  Name      string
  Path      string
  SharePath string
}

func NewWorker (lang *Lang, id int64, in chan Command) *Worker {
  wName := "perdoker_" + lang.Name +"_" + strconv.FormatInt(id, 10)
  path := "/tmp/" + lang.Name + "/"
  sharePath := path + ":" + path

  w := &Worker{lang, id, in, maxExecuteSeconds * time.Second, wName, path, sharePath }
  w.Start()
  return w
}

func (w *Worker) Start () {
  log.Println("Starting", w.Lang.Name, "worker ", w.Id)


  go func () {

    for {
      c := <- w.in
      log.Println("Worker", w.Id, ". Precessing", w.Lang.Name, "...")

      filePath := w.Path + w.Lang.uniqFileName() 

      ioutil.WriteFile(filePath, []byte(c.Command()), 755)

      //clear old container
      rm := exec.Command("docker", "rm", w.Name)
      rm.Start()
      rm.Wait()

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
        cmd.Process.Kill()
        //exec.Command("docker", "kill", w.Name).Run()
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
