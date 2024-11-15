package repository

import (
	"github.com/svnotn/test-port/port-service/internal/config"
	"github.com/svnotn/test-port/port-service/internal/domain"
	"github.com/svnotn/test-port/port-service/internal/model"
	"github.com/svnotn/test-port/port-service/internal/storage"
	"github.com/svnotn/test-port/port-service/internal/storage/in_mem"
)

type (
	PortRepository interface {
		Add(port model.Port) error
		Remove(port model.Port) error
		GetBy(port model.Port) (domain.Port, error)
	}

	PortRepo struct {
		s storage.Storage
	}
)

func New(config config.Port) PortRepository {
	r := &PortRepo{
		s: in_mem.New(config.CountIn, config.CountOut),
	}
	return r
}

func (r *PortRepo) Add(port model.Port) error {
	return r.s.Add(port)
}

func (r *PortRepo) Remove(port model.Port) error {
	return r.s.Remove(port)
}

func (r *PortRepo) GetBy(port model.Port) (domain.Port, error) {
	return r.s.GetBy(port)
}
