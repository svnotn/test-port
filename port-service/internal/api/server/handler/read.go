package handler

import (
	"errors"
	"strconv"

	"github.com/svnotn/test-port/port-service/internal/api/server/util"
	"github.com/svnotn/test-port/port-service/internal/model"
	"github.com/svnotn/test-port/port-service/internal/service/worker"

	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	labelRead = "[HTTP/API/READ ]: "
)

func NewReadHandler(worker *worker.Worker) func(*routing.Context) error {
	return func(c *routing.Context) error {
		var (
			id  int
			err error
		)
		params := util.GetQueryParams(c.QueryArgs())
		for i := range params {
			param, value := string(params[i].Key), string(params[i].Value)
			switch param {
			case "id":
				id, err = strconv.Atoi(value)
				if err != nil {
					log.Error(labelRead, err.Error())
					return errors.New("invalid id")
				}
			default:
				log.Warn(labelRead, "unknown parameter: ", param)
			}
		}
		cmd := model.NewCommand(id, model.Read, 0)
		worker.Send(cmd)
		if r := cmd.Result(); !r.Ok {
			c.Error(r.Err.Error(), fasthttp.StatusInternalServerError)
			return err
		}
		return nil
	}
}