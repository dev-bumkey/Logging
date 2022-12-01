package service

import (
	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/database"
	"github.com/cocktailcloud/acloud-alarm-collector/application/utils"
)

type Service interface {
}

type AlarmService struct {
	Context    *config.Context
	DB         database.DBAdapter
	MainQueue  utils.AlarmQueue
	RetryQueue utils.AlarmQueue
}

//func New(adapter database.DBAdapter, mq utils.AlarmQueue, rq utils.AlarmQueue) (AlarmService, error) {
func New(context *config.Context) (AlarmService, error) {
	service := AlarmService{
		Context:    context,
		DB:         context.DBAdapter,
		MainQueue:  context.MainQueue,
		RetryQueue: context.RetryQueue,
	}
	service.ReloadMemory()
	return service, nil
}
