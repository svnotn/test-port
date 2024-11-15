package worker

import (
	"context"
	"fmt"

	"github.com/svnotn/test-port/port-service/internal/config"
	"github.com/svnotn/test-port/port-service/internal/domain"
	"github.com/svnotn/test-port/port-service/internal/model"
	"github.com/svnotn/test-port/port-service/internal/repository"
	"github.com/svnotn/test-port/port-service/internal/storage"

	log "github.com/sirupsen/logrus"
)

const (
	label = "[WORKER]: \t%s"
)

type Worker struct {
	attempts int
	repo     storage.Storage
	ch       chan *model.Command
}

func New(config config.Worker, repo repository.PortRepository) *Worker {
	w := &Worker{
		attempts: config.AttemptsCount,
		repo:     repo,
		ch:       make(chan *model.Command, config.BuffSize),
	}
	return w
}

func (w *Worker) Send(cmd *model.Command) {
	w.ch <- cmd
}

func (w *Worker) Run(ctx context.Context) {
	log.Info(fmt.Sprintf(label, "start"))
	for {
		select {
		case cmd := <-w.ch:
			if err := execute(w.repo, cmd, w.attempts); err != nil {
				log.Error(fmt.Sprintf(label, err.Error()))
			}
		case <-ctx.Done():
			log.Info(fmt.Sprintf(label, "stop"))
			return
		}
	}
}

func execute(repo repository.PortRepository, cmd *model.Command, attempts int) error {
	port, err := repo.GetBy(cmd.ToPort())
	status := false

	defer func() {
		cmd.SetResult(model.Result{status, err})
	}()

	if err != nil {
		return err
	}

	for i := range attempts {
		err = port.Open()
		if port.State() != domain.Opened || err != nil {
			log.Error(fmt.Sprintf(label, fmt.Sprintf("can't open port {%d, %d}. retry %d", cmd.Action, cmd.ID, i)))
		} else {
			break
		}
	}
	if port.State() != domain.Opened {
		return err
	}

	defer func() {
		if err = port.Close(); err != nil {
			log.Error(fmt.Sprintf(label, err.Error()))
		}
	}()

	switch cmd.Action {
	case model.Read:
		if _, err = port.Read(); err != nil {
			return err
		}
	case model.Write:
		if err = port.Write(cmd.Transaction); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown action: %d", cmd.Action)
	}
	status = true
	return nil
}
