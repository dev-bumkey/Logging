package scheduler

import (
	"fmt"

	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/robfig/cron/v3"
)

type ReloadMemoryScheduler struct {
	cr      *cron.Cron
	Context *config.Context
	service service.AlarmService
}

func NewReloadMemoryJob(context *config.Context, alarmService service.AlarmService) (ReloadMemoryScheduler, error) {

	c := cron.New(cron.WithSeconds())

	interval := fmt.Sprintf("@every %ds", 60)

	logger.Debug("memory reload cron job setting: ", interval)
	if _, err := c.AddFunc(interval, alarmService.ReloadMemory); err != nil {
		return ReloadMemoryScheduler{}, fmt.Errorf("fail to set metric collecting cron job: %s", err.Error())
	} else {
		logger.Info("setup memory reload cron job: ok")
	}

	c.Start()

	return ReloadMemoryScheduler{
		cr:      c,
		Context: context,
		service: alarmService,
	}, nil
}

func (s ReloadMemoryScheduler) Run() error {
	s.cr.Start()
	return nil
}

func (s ReloadMemoryScheduler) Stop() error {
	if s.cr != nil {
		s.cr.Stop()
	}
	return nil
}
