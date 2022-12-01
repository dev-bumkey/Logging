package scheduler

import (
	"fmt"

	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/robfig/cron/v3"
)

type OrphanScheduler struct {
	cr      *cron.Cron
	Context *config.Context
}

func NewOrphanScheduler(context *config.Context, alarmService service.AlarmService) (OrphanScheduler, error) {
	c := cron.New(cron.WithSeconds())

	// interval := "0 0 12 * * *"
	interval := fmt.Sprintf("@every %ds", 60)
	logger.Debug("main-queue cron job setting: ", interval)
	if _, err := c.AddFunc(interval, alarmService.OrphanProcess); err != nil {
		return OrphanScheduler{}, fmt.Errorf("fail to set metric collecting cron job: %s", err.Error())
	} else {
		logger.Info("setup main-queue cron job: ok")
	}

	c.Start()

	return OrphanScheduler{
		cr:      c,
		Context: context,
	}, nil
}

func (s OrphanScheduler) Run() error {
	s.cr.Start()
	return nil
}

func (s OrphanScheduler) Stop() error {
	if s.cr != nil {
		s.cr.Stop()
	}
	return nil
}
