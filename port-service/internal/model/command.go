package model

type Action int

const (
	Read Action = iota
	Write
)

type Result struct {
	Ok  bool
	Err error
}

type Command struct {
	Action      Action
	ID          int
	Transaction int

	ch chan Result
}

func NewCommand(id int, action Action, transaction int) *Command {
	return &Command{
		Action:      action,
		ID:          id,
		Transaction: transaction,
		ch:          make(chan Result),
	}
}

func (c *Command) ToPort() Port {
	return Port{Type: PortType(c.Action), ID: c.ID}
}

func (c *Command) SetResult(r Result) {
	c.ch <- r
}

func (c *Command) Result() Result {
	r := <-c.ch
	close(c.ch)
	return r
}
