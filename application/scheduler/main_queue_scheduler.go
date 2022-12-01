package scheduler

import (
	"fmt"

	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/robfig/cron/v3"
)

// @Deprecated
type MainQueueScheduler struct {
	cr      *cron.Cron
	Config  *config.Config
	Service service.AlarmService
}

func NewMainQueueJob(conf *config.Config, alarmService service.AlarmService) (MainQueueScheduler, error) {

	c := cron.New(cron.WithSeconds())

	interval := fmt.Sprintf("@every %ds", 1)

	logger.Debug("main-queue cron job setting: ", interval)
	if _, err := c.AddFunc(interval, alarmService.MainQueueProcess); err != nil {
		return MainQueueScheduler{}, fmt.Errorf("fail to set metric collecting cron job: %s", err.Error())
	} else {
		logger.Info("setup main-queue cron job: ok")
	}

	c.Start()

	return MainQueueScheduler{
		cr:      c,
		Config:  conf,
		Service: alarmService,
	}, nil
}

func (s MainQueueScheduler) Run() error {
	s.cr.Start()
	return nil
}

func (s MainQueueScheduler) Stop() error {
	if s.cr != nil {
		s.cr.Stop()
	}
	return nil
}
