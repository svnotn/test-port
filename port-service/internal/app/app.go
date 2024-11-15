package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/svnotn/test-port/port-service/internal/api/server"
	"github.com/svnotn/test-port/port-service/internal/config"
	"github.com/svnotn/test-port/port-service/internal/model"
	"github.com/svnotn/test-port/port-service/internal/repository"
	"github.com/svnotn/test-port/port-service/internal/service/worker"

	log "github.com/sirupsen/logrus"
)

const label = "[APP]: \t"

type application struct {
	config config.Application
	repo   repository.PortRepository
	worker *worker.Worker
	server *server.Server
}

func newApp(config *config.Config) *application {
	app := &application{
		config: config.Application,
		repo:   repository.New(config.Port),
	}
	wrk := worker.New(config.Worker, app.repo)
	app.worker = wrk

	srv := server.New(config.Server, app.worker)
	app.server = srv

	return app
}

func (app *application) fillRep(count int, portType model.PortType) {
	for i := range count {
		err := app.repo.Add(model.Port{Type: portType, ID: i})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *application) Setup(config *config.Config) {
	app.fillRep(config.Port.CountIn, model.TypeIN)
	app.fillRep(config.Port.CountIn, model.TypeOUT)
}

func (app *application) Run(ctx context.Context) error {
	log.Info(label, "start")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		app.worker.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		app.server.Start(ctx)
	}()

	<-ctx.Done()

	if err := gracefulShutdown(wg, app.config.GracefulTimeout); err != nil {
		return err
	}

	log.Info(label, "stop")
	return nil
}

func Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(label, err)
	}
	cfg.Print()

	app := newApp(cfg)
	app.Setup(cfg)
	if err := app.Run(ctx); err != nil {
		log.Fatal(label, err)
	}
}

func gracefulShutdown(wg *sync.WaitGroup, waitTimeout time.Duration) error {
	wgDone := make(chan struct{})
	go func() {
		wg.Wait()
		wgDone <- struct{}{}
	}()

	select {
	case <-time.After(waitTimeout):
		return errors.New("failed to shutdown gracefully")
	case <-wgDone:
		return nil
	}
}
