package perd

// Command is an interface that must be implemented to in order to process
// commands.
type Command interface {
	Response([]byte, []byte, int)
	Command() string
	Language() *Lang
}

// NewCommand returns new Command
func NewCommand(lang *Lang, c string, r chan Result) Command {
	return &command{c, r, lang}
}

type command struct {
	command         string
	responseChannel chan Result
	lang            *Lang
}

func (c *command) Response(out, err []byte, code int) {
	c.responseChannel <- NewResult(out, err, code)
}

func (c *command) Command() string {
	return c.command
}

func (c *command) Language() *Lang {
	return c.lang
}
