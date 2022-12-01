package scheduler

import (
	"fmt"

	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/robfig/cron/v3"
)

type RetryScheduler struct {
	cr      *cron.Cron
	Context *config.Context
}

func NewRetryScheduler(context *config.Context, alarmService service.AlarmService) (RetryScheduler, error) {
	c := cron.New(cron.WithSeconds())

	interval := fmt.Sprintf("@every %ds", 5)
	logger.Debug("retry cron job setting: ", interval)
	if _, err := c.AddFunc(interval, alarmService.RetryProcess); err != nil {
		return RetryScheduler{}, fmt.Errorf("fail to retry cron job: %s", err.Error())
	} else {
		logger.Info("setup retry cron job: ok")
	}

	c.Start()

	return RetryScheduler{
		cr:      c,
		Context: context,
	}, nil
}

func (s RetryScheduler) Run() error {
	s.cr.Start()
	return nil
}

func (s RetryScheduler) Stop() error {
	if s.cr != nil {
		s.cr.Stop()
	}
	return nil
}
