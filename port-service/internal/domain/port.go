package domain

type State int

const (
	Opened State = 1 + iota
	Closed
)

type Port interface {
	Open() error
	Close() error
	State() State
	Read() (int, error)
	Write(transaction int) error
}
