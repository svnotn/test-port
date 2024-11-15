package portout

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/svnotn/test-port/port-service/internal/domain"
)

const (
	label = "[PORT_OUT]: \t%s(%d, %d)"
)

type Port struct {
	id          int
	transaction int
	state       domain.State

	mu      sync.Mutex
	chRead  chan int
	chClose chan bool
}

func New(id int) domain.Port {
	p := &Port{
		id:          id,
		transaction: -1,
		state:       domain.Closed,
		chRead:      make(chan int),
		chClose:     make(chan bool),
	}
	return p
}

func (p *Port) Open() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == domain.Opened {
		return fmt.Errorf("port out %d already opened", p.id)
	}
	p.state = domain.Opened
	go p.run()
	return nil
}

func (p *Port) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == domain.Closed {
		return fmt.Errorf("port out %d already closed", p.id)
	}
	p.state = domain.Closed
	p.chClose <- true
	return nil
}

func (p *Port) State() domain.State {
	return p.state
}

func (p *Port) Read() (int, error) {
	return 0, nil // port int does not support this method
}

func (p *Port) Write(transaction int) error {
	if p.state != domain.Opened {
		return fmt.Errorf("port out %d not opened", p.id)
	}
	p.chRead <- transaction
	return nil
}

func (p *Port) run() {
	for {
		select {
		case v := <-p.chRead:
			p.mu.Lock()
			p.transaction = v
			log.Info(fmt.Sprintf(label, "WRITE", p.id, v))
			p.mu.Unlock()
		case <-p.chClose:
			return
		}
	}
}
