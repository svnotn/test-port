package model

type PortType int

const (
	TypeIN = iota
	TypeOUT
)

type Port struct {
	Type PortType
	ID   int
}
