package perd

import "os/exec"
import "log"
import "io/ioutil"

type RubyRunner struct {
  runCh chan Command
}

func NewRubyRunner() Runner {
  return &RubyRunner{make(chan Command)}
}

func (r *RubyRunner) Run () {
  go func () {

    for {
      c := <- r.runCh

      ioutil.WriteFile("/tmp/1.rb", []byte(c.Command()), 755)
      exec.Command("docker", "rm", "ruby").Run()

      out, err := exec.Command("docker", "run", "-v", "/tmp:/tmp/host", "-name=ruby", "fd61e37b54de", "/bin/bash", "-l", "-c", "ruby /tmp/host/1.rb").CombinedOutput()

      log.Println("OUT: ", string(out), " Error:", err)
      c.Response(string(out))
    }

  }()
}

func (r *RubyRunner) Eval (command string) Result {
  respCh := make(chan Result)
  r.runCh <- NewCommand(command, respCh)
  return <-respCh
}


