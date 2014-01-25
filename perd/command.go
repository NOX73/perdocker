package perd

// Command is an interface that must be implemented to in order to process
// commands.
type Command interface {
	Response([]byte, []byte, int)
	Command() string
}

// NewCommand returns new Command
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
