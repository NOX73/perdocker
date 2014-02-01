package perd

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	eol byte = 10
)

var (
	// ErrCantStart indicates fail in starting particular container (detailed error
	// will be printed to the stdout.
	ErrCantStart = errors.New("can't start container")
)
var Backend BackendI = new(backend)

type BackendI interface {
	Start(name, image, shared, mem, cpu string) (inCh, outCh, errCh chan []byte, err error)
	Stop(name string)
	DetectForks(name string) (found bool, err error)
}

type backend struct{}

func (b *backend) Start(name, image, shared, mem, cpu string) (inCh, outCh, errCh chan []byte, err error) {
	cmd := exec.Command("docker", "run", "-m", mem, "-c", cpu, "-i", "-v", shared, "-name="+name, image, "/bin/bash", "-l")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	err = b.waitStart(name)
	if err != nil {
		io.Copy(os.Stdout, stderr)

		stderr.Close()
		stdout.Close()
		stdin.Close()

		return
	}

	inCh = make(chan []byte, 5)
	outCh = make(chan []byte, 5)
	errCh = make(chan []byte, 5)

	go copyStreams(stdin, stdout, stderr, inCh, outCh, errCh)

	return
}

func (b *backend) Stop(name string) {
	for b.isExist(name) {
		b.kill(name)
		b.rm(name)
	}
}

// DetectForks returns true if number of processes greater than 1
// This may be a result of the fork bomb.
func (c *backend) DetectForks(name string) (found bool, err error) {
	c1 := exec.Command("docker", "top", name)
	c2 := exec.Command("wc", "-l")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()

	nprocStr := b2.String()
	var nproc int
	nproc, err = strconv.Atoi(nprocStr[:len(nprocStr)-1])
	if err != nil {
		return false, err
	}

	found = (nproc - 1) > 1
	return found, nil
}

func (b *backend) waitStart(name string) error {
	for i := 0; i < 5; i++ {
		if b.isExist(name) {
			return nil
		}
	}
	return ErrCantStart
}

func (b *backend) isExist(name string) bool {
	err := exec.Command("docker", "inspect", name).Run()
	return err == nil
}

func (b *backend) rm(name string) error {
	return exec.Command("docker", "rm", name).Run()
}

func (b *backend) kill(name string) error {
	return exec.Command("docker", "kill", name).Run()
}

func copyStreams(in io.WriteCloser, out, err io.ReadCloser, inCh, outCh, errCh chan []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Copy stream finished with panic(). Reason: ", r)
		} else {
			log.Println("Copy stream finished correctly.")
		}
	}()

	go readLinesToChannel(out, outCh)
	go readLinesToChannel(err, errCh)

	writeLinesFromChannel(in, inCh)

	in.Close()
	out.Close()
	err.Close()

	close(outCh)
	close(errCh)

	return
}

func readLinesToChannel(rc io.ReadCloser, ch chan []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Read stream finished with panic(). Reason: ", r)
		} else {
			log.Println("Read stream finished correctly.")
		}
	}()

	r := bufio.NewReader(rc)

	for {
		line, err := r.ReadBytes(eol)
		if err != nil {
			break
		}
		ch <- line
	}

}

func writeLinesFromChannel(wc io.WriteCloser, ch chan []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Write stream finished with panic(). Reason: ", r)
		} else {
			log.Println("Write stream finished correctly.")
		}
	}()

	w := bufio.NewWriter(wc)

	for {
		line, ok := <-ch

		if !ok {
			break
		}

		_, err := w.Write(line)
		if err != nil {
			break
		}

		err = w.WriteByte(eol)
		if err != nil {
			break
		}

		w.Flush()
	}
}
