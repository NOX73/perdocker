package perd

type Command interface {
  Response(string)
  Command()string
}

func NewCommand(c string, r chan Result) Command {
  return &command{c, r}
}

type command struct {
  command string
  responseChannel chan Result
}

func (c *command) Response (r string) {
  c.responseChannel <- NewResult(r)
}

func (c *command) Command () string {
  return c.command
}


