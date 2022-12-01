package api

import (
	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
)

type API struct {
	Config  *config.Config
	Service service.AlarmService
}

func New(conf *config.Config, service service.AlarmService) (*API, error) {
	api := &API{
		Config:  conf,
		Service: service,
	}
	return api, nil
}
