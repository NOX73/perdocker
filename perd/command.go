package perd

type Command interface {
	Response([]byte, []byte, int)
	Command() string
}

func NewCommand(c string, r chan Result) Command {
	return &command{c, r}
}

type command struct {
	command         string
	responseChannel chan Result
}

func (c *command) Response(out, err []byte, code int) {
	c.responseChannel <- NewResult(out, err, code)
}

func (c *command) Command() string {
	return c.command
}
