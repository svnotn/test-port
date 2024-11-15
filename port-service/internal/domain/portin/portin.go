package portin

import (
	"fmt"
	"math/rand"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/svnotn/test-port/port-service/internal/domain"
)

const (
	label = "[PORT_IN ]: \t%s(%d, %d)"
)

type Port struct {
	id    int
	value int
	state domain.State

	mu      sync.Mutex
	chRead  chan int
	chClose chan bool
}

func New(id int) domain.Port {
	p := &Port{
		id:      id,
		value:   -1,
		state:   domain.Closed,
		chRead:  make(chan int),
		chClose: make(chan bool),
	}
	return p
}

func (p *Port) Open() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == domain.Opened {
		return fmt.Errorf("port in %d already opened", p.id)
	}
	p.state = domain.Opened
	go p.run()
	return nil
}

func (p *Port) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == domain.Closed {
		return fmt.Errorf("port in %d already closed", p.id)
	}
	p.state = domain.Closed
	p.chClose <- true
	return nil
}

func (p *Port) State() domain.State {
	return p.state
}

func (p *Port) Read() (int, error) {
	if p.state != domain.Opened {
		return 0, fmt.Errorf("port in %d not opened", p.id)
	}
	v := rand.Intn(2)
	p.chRead <- v
	return v, nil
}

func (p *Port) Write(_ int) error {
	return nil // port int does not support this method
}

func (p *Port) run() {
	for {
		select {
		case v := <-p.chRead:
			p.mu.Lock()
			p.value = v
			log.Info(fmt.Sprintf(label, "READ", p.id, p.value))
			p.mu.Unlock()
		case <-p.chClose:
			return
		}
	}
}
