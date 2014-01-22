package perd

import (
  "io"
  "os/exec"
  "strconv"
  "bufio"
)

type Container interface {
  Start()
  Stop()
  Restart()

  StdIn()   chan<- []byte
  StdOut()  <-chan []byte
  StdErr()  <-chan []byte
}

const (
  MemLimit = "10m"
  CpuLimit = "1"
)

type container struct {
  Id      int64
  Lang    *Lang

  cmd     *exec.Cmd

  name string

  tmpHost string
  tmpGuest string

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

  inWriter *bufio.Writer
  outReader *bufio.Reader
  errReader *bufio.Reader

  inCh  chan []byte
  outCh  chan []byte
  errCh  chan []byte
}

func NewContainer(id int64, lang *Lang) Container {

  name := "perdoker_" + lang.Name + "_" + strconv.FormatInt(id, 10)

  c := &container{
    Id: id,
    Lang: lang,

    name: name,
    tmpHost: "/tmp/perdocker/" + lang.Name + "/" + name + "/",
    tmpGuest: "/tmp/perdocker/",
  }

  return c
}

func (c *container) Start () {
  cmd := exec.Command("docker", "run", "-m", MemLimit, "-c", CpuLimit, "-i", "-v", c.sharedPaths(), "-name=" + c.name, c.Lang.Image, "/bin/bash", "-l")
  c.cmd = cmd

	c.stdin, _ = cmd.StdinPipe()
	c.stdout, _ = cmd.StdoutPipe()
	c.stderr, _ = cmd.StderrPipe()

	c.inWriter = bufio.NewWriter(c.stdin)
	c.outReader = bufio.NewReader(c.stdout)
	c.errReader = bufio.NewReader(c.stderr)

  go readLinesToChannel(c.outReader, c.outCh)
  go readLinesToChannel(c.errReader, c.errCh)

}

func readLinesToChannel(r *bufio.Reader, ch chan []byte) {
  for {
    line, err := r.ReadBytes(eol)
    if err != nil { break }
    ch <- line
  }
}


func (c *container) sharedPaths () string {
  return c.tmpHost + ":" + c.tmpGuest + ":ro"
}

func (c *container) Stop () {
	c.stdin.Close()
	c.stdout.Close()
	c.stderr.Close()

	c.clear()
}

func (c *container) Restart () {
  c.Stop()
  c.Start()
}

func (c *container) StdIn () chan<- []byte {
  return c.inCh
}

func (c *container) StdOut () <-chan []byte {
  return c.outCh
}

func (c *container) StdErr () <-chan []byte {
  return c.errCh
}

func (c *container) clear () {
	for c.isExist() {
		c.kill()
		c.rm()
	}
}

func (c *container) rm () error {
	return exec.Command("docker", "rm", c.name).Run()
}

func (c *container) kill () error {
	return exec.Command("docker", "kill", c.name).Run()
}

func (c *container) isExist () bool {
	err := exec.Command("docker", "inspect", c.name).Run()
	return err == nil
}
