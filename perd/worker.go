package perd

import (
  "strconv"
  "io/ioutil"
  "os/exec"
  "bytes"
  "syscall"
  "log"
)

type Worker struct {
  Lang      *Lang
  Id        int64
  in        chan Command
}

func NewWorker (lang *Lang, id int64, in chan Command) *Worker {
  w := &Worker{lang, id, in}
  w.Start()
  return w
}

func (w *Worker) Start () {
  path := "/tmp/" + w.Lang.Name + "/"
  wName := "perdoker_" + w.Lang.Name +"_" + strconv.FormatInt(w.Id, 10)
  sharePath := path + ":" + path

  go func () {

    for {
      c := <- w.in

      filePath := path + w.Lang.uniqFileName() 

      ioutil.WriteFile(filePath, []byte(c.Command()), 755)
      exec.Command("docker", "rm", wName).Run()

      cmd := exec.Command("docker", "run", "-v", sharePath, "-name=" + wName, w.Lang.Image, "/bin/bash", "-l", "-c", w.Lang.RunCommand(filePath))

      var stdOut, stdErr bytes.Buffer
      var code int
      cmd.Stdout, cmd.Stderr = &stdOut, &stdErr

      cmd.Start()
      err := cmd.Wait()

      code = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

      if err != nil { log.Println("Error:", err) }
      c.Response(stdOut.Bytes(), stdErr.Bytes(), code)
    }

  }()

}
