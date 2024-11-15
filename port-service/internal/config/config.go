package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type (
	Env struct {
		Path string
	}

	Config struct {
		Application Application
		Server      Server
		Port        Port
		Worker      Worker
	}

	Application struct {
		CpuLimit        uint32        `env:"APPLICATION_CPU_LIMIT"        long:"cpu-limit" envDefault:"1"`
		GracefulTimeout time.Duration `env:"APPLICATION_GRACEFUL_TIMEOUT"                  envDefault:"30s"`
	}

	Server struct {
		Port    int           `env:"SERVER_PORT"    envDefault:"8122"`
		Timeout time.Duration `env:"SERVER_TIMEOUT" envDefault:"5s"`
	}

	Port struct {
		CountIn  int `env:"COUNT_IN"  envDefault:"1" validate:"required,min=1"`
		CountOut int `env:"COUNT_OUT" envDefault:"1" validate:"required,min=1"`
	}

	Worker struct {
		BuffSize      int `env:"WORKER_BUFF_SIZE"      envDefault:"1" validate:"required,min=1"`
		AttemptsCount int `env:"WORKER_ATTEMPTS_COUNT" envDefault:"1" validate:"required,min=1"`
	}
)

func newEnv() *Env {
	e := &Env{}
	flag.StringVar(&e.Path, "path", "", "path to env file. must be set")
	flag.Parse()
	if e.Path == "" {
		log.Fatal("[CONF]: ", fmt.Errorf("path is required"))
	}
	return e
}

func (e *Env) setup() error {
	err := godotenv.Load(e.Path)
	if err != nil {
		return err
	}
	return nil
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cfg.parse(); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) Print() {
	log.Info("============= [CONFIG] =============")
	log.Info("[APP]    CpuLimit.......: ", cfg.Application.CpuLimit)
	log.Info("[APP]    GracefulTimeout: ", cfg.Application.GracefulTimeout)
	log.Info("[SRV]    Port...........: ", cfg.Server.Port)
	log.Info("[SRV]    Timeout........: ", cfg.Server.Timeout)
	log.Info("[PORT]   IN.............: ", cfg.Port.CountIn)
	log.Info("[PORT]   OUT............: ", cfg.Port.CountOut)
	log.Info("[WORKER] BuffSize.......: ", cfg.Worker.BuffSize)
	log.Info("[WORKER] Attempts.......: ", cfg.Worker.AttemptsCount)
	log.Info("====================================")
}

func (cfg *Config) parse() error {
	e := newEnv()
	if err := e.setup(); err != nil {
		return err
	}
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) validate() error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return err
	}
	return nil
}
