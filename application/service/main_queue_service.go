package service

import (
	"time"

	"github.com/cocktailcloud/acloud-alarm-collector/application/model"
	"github.com/cocktailcloud/acloud-alarm-collector/application/structure"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
)

const (
	updateAlarmTimestamp = "update alarm_current set timestamp = now() AT TIME ZONE 'UTC' where cluster_id=$1 and alarm_id=$2"
	deleteAlarm          = "delete from alarm_current where cluster_id=$1 and alarm_id=$2"
)

func (s *AlarmService) LoopMainQueueProcess() {
	for {
		s.MainQueueProcess()
	}
}

func (s *AlarmService) MainQueueProcess() {
	if s.MainQueue.IsEmpty() {
		time.Sleep(1 * time.Second)
		// logger.Debug("Retry Queue is empty")
		return
	}
	contents, ok := s.MainQueue.Poll()
	// logger.Debugf("retry entry is %d'th: %v", contents.Retry, string(contents.Contents))
	if ok {
		err := s.process(contents)
		if err != nil {
			logger.Errorf("fail to retry %v", err)
			contents.Retry += 1
			s.RetryQueue.Add(contents)
		}

	}
}

// Filtering
// Exists Check
// Action (Insert, Update, Delete)
// Post Action
func (s *AlarmService) process(alarm structure.AlarmItem) error {

	// TODO:: Filtering
	// return fmt.Errorf("test Error")

	// Exists check
	existsAlarm := s.existsAlarm(alarm.Alarm)
	isFired := isFiring(alarm.Alarm)
	var err error
	var mode string = "NONE"

	if isFired {
		if existsAlarm {
			mode = "UPDATE"
			err = s.updateAlarm(alarm.Alarm)
			checkErr(err, "")
		} else {
			mode = "INSERT"
			err = s.insertAlarm(alarm.Alarm)
			checkErr(err, "")
		}
	} else {
		if existsAlarm {
			mode = "CLEAR"
			err = s.clearAlarm(alarm.Alarm)
			checkErr(err, "")
		}
	}

	logger.Infof("execute alarm[%s]: %s, %v, %v, %v", mode, alarm.Alarm.ClusterId, alarm.Alarm.Alertname, alarm.Alarm.Status, alarm.Alarm.Fingerprint)

	return err
}

func (s *AlarmService) existsAlarm(alarm *structure.Alarm) bool {
	copiedList := s.Context.GetAliveAlarm()
	if items, ok := copiedList[alarm.ClusterId]; ok {
		if _, has := items[alarm.Fingerprint]; has {
			return true
		}
	}
	return false
}

func (s *AlarmService) updateAlarm(alarm *structure.Alarm) error {
	_, err := s.DB.DB.GetClient().Exec(updateAlarmTimestamp, alarm.ClusterId, alarm.Fingerprint)

	return err
}

func (s *AlarmService) insertAlarm(alarm *structure.Alarm) error {
	// txDB, _ := s.DB.DB.BeginTransaction()
	currentModel := &model.AlarmCurrent{
		ClusterId:    alarm.ClusterId,
		AlarmId:      alarm.Fingerprint,
		GroupLabel:   alarm.GroupKey,
		Namespace:    alarm.Namespace,
		WorkloadType: alarm.WorkloadType,
		WorkloadName: alarm.Workload,
		Pod:          alarm.Pod,
		Container:    alarm.Container,
		AlarmName:    alarm.Alertname,
		Severity:     alarm.Severity,
		Description:  alarm.Description,
		Status:       alarm.Status,
		Startsat:     alarm.StartsAt,
		AlarmKey:     alarm.AlarmKey,
		Endsat:       alarm.EndsAt,
		Timestamp:    alarm.Timestamp,
	}
	err := s.DB.DB.GetClient().Insert(currentModel)
	// 메모리 작업
	s.UpdateMemoryNew(alarm.ClusterId, alarm.Fingerprint)

	return err
}

func (s *AlarmService) clearAlarm(alarm *structure.Alarm) error {
	txDB, _ := s.DB.DB.BeginTransaction()
	var err error
	historyModel := &model.AlarmHistory{
		ClusterId:    alarm.ClusterId,
		AlarmId:      alarm.Fingerprint,
		Namespace:    alarm.Namespace,
		WorkloadType: alarm.WorkloadType,
		WorkloadName: alarm.Workload,
		Pod:          alarm.Pod,
		Container:    alarm.Container,
		AlarmName:    alarm.Alertname,
		Severity:     alarm.Severity,
		Description:  alarm.Description,
		Status:       alarm.Status,
		Startsat:     alarm.StartsAt,
		AlarmKey:     alarm.AlarmKey,
		Endsat:       alarm.EndsAt,
		Timestamp:    alarm.Timestamp,
	}

	// Delete Alarm
	_, err = txDB.GetClient().Exec(deleteAlarm, alarm.ClusterId, alarm.Fingerprint)
	if err != nil {
		txDB.Rollback()
		return err
	}

	// Move to History
	err = txDB.GetClient().Insert(historyModel)
	if err != nil {
		txDB.Rollback()
		return err
	}

	txDB.Commit()
	// 메모리 작업
	s.UpdateMemoryClear(alarm.ClusterId, alarm.Fingerprint)
	return nil
}

func isFiring(alarm *structure.Alarm) bool {
	return alarm.Status == "firing"
}
