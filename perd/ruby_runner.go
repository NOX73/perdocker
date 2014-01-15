package perd

import "os/exec"
import "log"
import "io/ioutil"
import "strconv"
import "bytes"
import "syscall"

const (
  path = "/tmp/ruby/"
  image = "perdocker/ruby"
)

type RubyRunner struct {
  runCh chan Command
  *runner
}

func NewRubyRunner() Runner {
  return &RubyRunner{make(chan Command), &runner{}}
}

var workerId int64
func (r *RubyRunner) RunWorker () {
  workerId ++

  wId := workerId
  wName := "perdoker_ruby_" + strconv.FormatInt(wId, 10)
  sharePath := path + ":" + path

  go func () {

    for {
      c := <- r.runCh

      filePath := path + uniqFileName() + ".rb"

      ioutil.WriteFile(filePath, []byte(c.Command()), 755)
      exec.Command("docker", "rm", wName).Run()

      cmd := exec.Command("docker", "run", "-v", sharePath, "-name=" + wName, image, "/bin/bash", "-l", "-c", "ruby " + filePath)

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

func (r *RubyRunner) RunWorkers (count int) {
  for i := count; i>0; i-- {
    r.RunWorker()
  }
}

func (r *RubyRunner) Eval (command string) Result {
  respCh := make(chan Result)
  r.runCh <- NewCommand(command, respCh)
  return <-respCh
}


