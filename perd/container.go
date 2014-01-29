package perd

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Container represents Docker container.
type Container interface {
	Init() error
	Clear() error
	Start() error
	Stop()
	Restart() error
	Exec(Command) (*Exec, error)
	Remove()
}

const (
	// MemLimit sets allowed memory limit
	MemLimit = "20m"

	// CPULimit sets allowed CPU count
	CPULimit = "1"
)

var (
	// ErrCantStart indicates fail in starting particular container (detailed error
	// will be printed to the stdout.
	ErrCantStart = errors.New("can't start container")
)

type container struct {
	ID   int64
	Lang *Lang

	cmd *exec.Cmd

	name string

	end []byte

	tmpHost  string
	tmpGuest string

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

// NewContainer returns new Container
func NewContainer(id int64, lang *Lang) (Container, error) {

	name := "perdoker_" + strconv.FormatInt(id, 10)
	tmpHost := "/tmp/perdocker/" + name + "/"
	tmpGuest := "/tmp/perdocker/"

	c := &container{
		ID:   id,
		Lang: lang,

		name: name,

		tmpHost:  tmpHost,
		tmpGuest: tmpGuest,
	}

	err := os.MkdirAll(c.tmpHost, 0755)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *container) echoEnd() error {
	var err error
	_, err = c.inWriter.WriteString("echo \"\n" + string(c.end) + "$?\"\n")
	_, err = c.inWriter.WriteString("echo " + string(c.end) + " 1>&2\n")
	c.inWriter.Flush()
	return err
}

func (c *container) Exec(command Command) (*Exec, error) {

	lang := command.Language()
	if lang == nil {
		lang = c.Lang
	}
	code := []byte(command.Command())

	var err error

	fileHost := c.tmpHost + lang.ExecutableFile()
	fileGuest := c.tmpGuest + lang.ExecutableFile()

	err = ioutil.WriteFile(fileHost, code, 0755)
	if err != nil {
		return nil, err
	}

	in := c.inWriter

	execStr := lang.RunCommand(fileGuest)

	_, err = in.WriteString(execStr + " 3>&- \n")
	err = c.echoEnd()
	if err != nil {
		return nil, err
	}

	exec := NewExec(c.outCh, c.errCh, c.end)

	return exec, nil
}

func (c *container) Start() error {
	var err error

	cmd := exec.Command("docker", "run", "-m", MemLimit, "-c", CPULimit, "-i", "-v", c.sharedPaths(), "-name="+c.name, c.Lang.Image, "/bin/bash", "-l")
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

func (c *container) Restart() error {
	c.Stop()
	return c.Start()
}

func (c *container) Init() error {
	var err error
	c.Remove()

	err = c.Start()
	if err != nil {
		return err
	}

	err = c.Clear()
	if err != nil {
		return err
	}

	return nil
}

func (c *container) Remove() {
	for c.isExist() {
		c.kill()
		c.rm()
	}
}

func (c *container) Clear() error {
	c.end = generateEnd()
	//TODO: Fork detector
	return c.clearStd()
}

func (c *container) clearStd() error {
	err := c.echoEnd()
	if err != nil {
		return err
	}
	exec := NewExec(c.outCh, c.errCh, c.end)

	return exec.Wait(5 * time.Second)
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

var endChars []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz0123456789")
var endLen = 30

func generateEnd() []byte {
	end := make([]byte, endLen)

	for i := 0; i < endLen; i++ {
		r := rand.Intn(len(endChars))
		end[i] = endChars[r]
	}

	return end
}

func (c *container) sharedPaths() string {
	return c.tmpHost + ":" + c.tmpGuest + ":ro"
}
