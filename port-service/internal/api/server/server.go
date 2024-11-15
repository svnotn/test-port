package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/svnotn/test-port/port-service/internal/api/server/handler"
	"github.com/svnotn/test-port/port-service/internal/config"
	"github.com/svnotn/test-port/port-service/internal/service/worker"

	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	label = "[SERVER]: \t%s"
)

type Server struct {
	config     config.Server
	httpServer *fasthttp.Server
	worker     *worker.Worker
}

func New(config config.Server, worker *worker.Worker) *Server {
	s := &Server{
		config: config,
		httpServer: &fasthttp.Server{
			Name: "port-service",
		},
		worker: worker,
	}
	return s
}

func (s *Server) Start(ctx context.Context) {
	s.initHandlers()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		log.Info(fmt.Sprintf(label, "starting on port: "), s.config.Port)

		err := s.httpServer.ListenAndServe(":" + strconv.Itoa(s.config.Port))
		if err != nil {
			log.Fatal(fmt.Sprintf(label, err))
		}
	}()

	<-ctx.Done()

	if err := s.httpServer.Shutdown(); err != nil {
		log.Error(fmt.Sprintf(label, err))
	}
	log.Info(fmt.Sprintf(label, "stop"))
}

func (s *Server) initHandlers() {
	router := routing.New()
	api := router.Group("/api")
	api.Post("/read", handler.NewReadHandler(s.worker))
	api.Post("/write", handler.NewWriteHandler(s.worker))
	s.httpServer.Handler = router.HandleRequest
}
