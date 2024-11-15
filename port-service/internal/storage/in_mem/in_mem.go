package in_mem

import (
	"fmt"

	"github.com/svnotn/test-port/port-service/internal/domain"
	"github.com/svnotn/test-port/port-service/internal/domain/portin"
	"github.com/svnotn/test-port/port-service/internal/domain/portout"
	"github.com/svnotn/test-port/port-service/internal/model"
	"github.com/svnotn/test-port/port-service/internal/storage"
)

type Storage struct {
	in  map[int]domain.Port
	out map[int]domain.Port
}

func New(in, out int) storage.Storage {
	s := &Storage{
		in:  make(map[int]domain.Port, in),
		out: make(map[int]domain.Port, out),
	}
	return s
}

func (s *Storage) Add(port model.Port) error {
	m := s.selectByType(port.Type)
	if _, ok := m[port.ID]; ok {
		return fmt.Errorf("port type %d with id %d already exists", port.Type, port.ID)
	}
	m[port.ID] = s.createByType(port)
	return nil
}

func (s *Storage) Remove(port model.Port) error {
	m := s.selectByType(port.Type)
	if _, ok := m[port.ID]; !ok {
		return fmt.Errorf("port type %d with id %d not found", port.Type, port.ID)
	}
	delete(m, port.ID)
	return nil
}

func (s *Storage) GetBy(port model.Port) (domain.Port, error) {
	m := s.selectByType(port.Type)
	if _, ok := m[port.ID]; !ok {
		return nil, fmt.Errorf("port type %d with id %d not found", port.Type, port.ID)
	}
	return m[port.ID], nil
}

func (s *Storage) selectByType(portType model.PortType) map[int]domain.Port {
	if portType == model.TypeIN {
		return s.in
	}
	return s.out
}

func (s *Storage) createByType(port model.Port) domain.Port {
	if port.Type == model.TypeIN {
		return portin.New(port.ID)
	}
	return portout.New(port.ID)
}
