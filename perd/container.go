package perd

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

type Container interface {
	Init()
	Clear()
	Start() error
	Stop()
	Restart()
	Exec([]byte) (*Exec, error)
	Remove()
}

const (
	MemLimit    = "10m"
	CpuLimit    = "1"
	WaitStarSec = 5
)

var (
	ErrCantStart = errors.New("Cant't start container.")
)

type container struct {
	Id   int64
	Lang *Lang

	cmd *exec.Cmd

	name string

	end     []byte
	command string

	tmpHost  string
	tmpGuest string

	fileHost  string
	fileGuest string

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	inWriter  *bufio.Writer
	outReader *bufio.Reader
	errReader *bufio.Reader

	inCh  chan []byte
	outCh chan []byte
	errCh chan []byte
}

func NewContainer(id int64, lang *Lang) (Container, error) {

	name := "perdoker_" + lang.Name + "_" + strconv.FormatInt(id, 10)
	tmpHost := "/tmp/perdocker/" + lang.Name + "/" + name + "/"
	tmpGuest := "/tmp/perdocker/"
	fileGuest := tmpGuest + lang.ExecutableFile()

	c := &container{
		Id:   id,
		Lang: lang,

		name: name,

		fileHost:  tmpHost + lang.ExecutableFile(),
		fileGuest: fileGuest,

		tmpHost:  tmpHost,
		tmpGuest: tmpGuest,

		command: lang.RunCommand(fileGuest),

		end: generateEnd(),
	}

	err := os.MkdirAll(c.tmpHost, 0755)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *container) Exec(file []byte) (*Exec, error) {
	var err error

	err = ioutil.WriteFile(c.fileHost, file, 0755)
	if err != nil {
		return nil, err
	}

	in := c.inWriter

	_, err = in.WriteString(c.command + " 3>&- \n")
	if err != nil {
		return nil, err
	}
	_, err = in.WriteString("echo " + string(c.end) + "$?\n")
	if err != nil {
		return nil, err
	}
	_, err = in.WriteString("echo " + string(c.end) + " 1>&2\n")
	if err != nil {
		return nil, err
	}

	in.Flush()

	exec := NewExec(c.outCh, c.errCh, c.end)

	go exec.Start()

	return exec, nil
}

func (c *container) Start() error {
	var err error

	cmd := exec.Command("docker", "run", "-m", MemLimit, "-c", CpuLimit, "-i", "-v", c.sharedPaths(), "-name="+c.name, c.Lang.Image, "/bin/bash", "-l")
	c.cmd = cmd

	c.stdin, _ = cmd.StdinPipe()
	c.stdout, _ = cmd.StdoutPipe()
	c.stderr, _ = cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = c.waitStart()
	if err != nil {
		c.rm()
		return err
	}

	c.inWriter = bufio.NewWriter(c.stdin)
	c.outReader = bufio.NewReader(c.stdout)
	c.errReader = bufio.NewReader(c.stderr)

	c.inCh = make(chan []byte, 5)
	c.outCh = make(chan []byte, 5)
	c.errCh = make(chan []byte, 5)

	go readLinesToChannel(c.outReader, c.outCh)
	go readLinesToChannel(c.errReader, c.errCh)

	return nil
}

func (c *container) Stop() {
	c.stdin.Close()
	c.stdout.Close()
	c.stderr.Close()

	close(c.inCh)
	close(c.outCh)
	close(c.errCh)

	c.Remove()
}

func (c *container) Restart() {
	c.Stop()
	c.Start()
}

func (c *container) Init() {
	c.Remove()
	c.Start()
}

func (c *container) Remove() {
	for c.isExist() {
		c.kill()
		c.rm()
	}
}

func (c *container) Clear() {
	//TODO: Fork detector
	//TODO: Clear stdOut stdErr
	//TODO: Generate end
}

func (c *container) rm() error {
	return exec.Command("docker", "rm", c.name).Run()
}

func (c *container) kill() error {
	return exec.Command("docker", "kill", c.name).Run()
}

func (c *container) isExist() bool {
	err := exec.Command("docker", "inspect", c.name).Run()
	return err == nil
}

func (c *container) waitStart() error {
	for i := 0; i < 5; i++ {
		if c.isExist() {
			return nil
		}
	}
	return ErrCantStart
}

func readLinesToChannel(r *bufio.Reader, ch chan []byte) {
	defer func() { recover() }()
	for {
		line, err := r.ReadBytes(eol)
		if err != nil {
			break
		}
		ch <- line
	}
}

func generateEnd() []byte {
	return []byte("asdfasdfasdfasdfgfsdfbewrv")
}

func (c *container) sharedPaths() string {
	return c.tmpHost + ":" + c.tmpGuest + ":ro"
}
