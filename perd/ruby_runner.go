package perd

import "os/exec"
import "log"
import "io/ioutil"
import "strconv"

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

      out, err := exec.Command("docker", "run", "-v", sharePath, "-name=" + wName, image, "/bin/bash", "-l", "-c", "ruby " + filePath).CombinedOutput()

      if err != nil { log.Println("Error:", err) }
      c.Response(out, []byte{}, 0)
    }

  }()
}

func (r *RubyRunner) Eval (command string) Result {
  respCh := make(chan Result)
  r.runCh <- NewCommand(command, respCh)
  return <-respCh
}


