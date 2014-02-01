package perd

import (
	"errors"
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

	Start() error
	Stop()
	Restart() error

	Exec(Command) (*Exec, error)

	Clear() error
}

const (
	// MemLimit sets allowed memory limit
	MemLimit = "20m"

	// CPULimit sets allowed CPU count
	CPULimit = "1"
)

var (
	ErrSendCommandTimeout = errors.New("Sendcommand to container timeout.")
)

type container struct {
	ID   int64
	Lang *Lang

	cmd *exec.Cmd

	name string

	end []byte

	tmpHost  string
	tmpGuest string

	inCh  chan<- []byte
	outCh <-chan []byte
	errCh <-chan []byte
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

func (c *container) Init() error {
	var err error

	err = c.Restart()
	if err != nil {
		return err
	}

	err = c.Clear()
	if err != nil {
		return err
	}

	return nil
}

func (c *container) Start() error {
	var err error
	c.inCh, c.outCh, c.errCh, err = Backend.Start(c.name, c.Lang.Image, c.sharedPaths(), MemLimit, CPULimit)

	return err
}

func (c *container) Stop() {
	Backend.Stop(c.name)
	close(c.inCh)
}

func (c *container) Restart() error {
	c.Stop()
	return c.Start()
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

	in := c.inCh

	execStr := lang.RunCommand(fileGuest)

	err = c.sendCommand(execStr + " 3>&- ")
	if err != nil {
		return nil, err
	}

	err = c.echoEnd()
	if err != nil {
		return nil, err
	}

	exec := NewExec(c.outCh, c.errCh, c.end)

	return exec, nil
}
func (c *container) Clear() error {
	c.end = generateEnd()
	//TODO: Fork detector
	return c.clearStd()
}

// Private

func (c *container) clearStd() error {
	err := c.echoEnd()
	if err != nil {
		return err
	}
	exec := NewExec(c.outCh, c.errCh, c.end)

	return exec.Wait(5 * time.Second)
}

func (c *container) sendCommand(command string) error {
	select {
	case c.inCh <- []byte(command):
	case <-time.After(5 * time.Second):
		return ErrSendCommandTimeout
	}
	return nil
}

func (c *container) echoEnd() error {
	var err error

	err = c.sendCommand("echo \"\n" + string(c.end) + "$?\"")
	if err != nil {
		return err
	}

	err = c.sendCommand("echo " + string(c.end) + " 1>&2")
	if err != nil {
		return err
	}

	return nil
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
