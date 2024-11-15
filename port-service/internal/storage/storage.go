package storage

import (
	"github.com/svnotn/test-port/port-service/internal/domain"
	"github.com/svnotn/test-port/port-service/internal/model"
)

type Storage interface {
	Add(port model.Port) error
	Remove(port model.Port) error
	GetBy(port model.Port) (domain.Port, error)
}
