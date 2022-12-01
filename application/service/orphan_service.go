package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cocktailcloud/acloud-alarm-collector/application/model"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
)

const (
	getOrphanAlarms    = "select alarm_id, cluster_id, timestamp from alarm_current where timestamp < now() - interval '%s'"
	removeOrphanAlarms = "delete from alarm_current where alarm_id in (%s)"
	moveOrphanAlarms   = `
	insert into alarm_history
	select cluster_id, alarm_id, namespace, workload_type, workload_name, pod, container, alarm_key, alarm_name, severity, description, 'resolved', startsat, endsat, 0, 'orphan alarm removed', null, null, null, now() AT TIME ZONE 'UTC'
	from alarm_current where alarm_id in (%s)`
)

// 기준이 될 기간 이상 업데이트 되지 않는 알람은 제거합니다. (고아로 판단)
func (s *AlarmService) OrphanProcess() {

	// var alarms []structure.AliveAlarm
	var alarms []model.AlarmCurrent
	exec := s.DB.DB.GetClient()
	_, err := exec.Select(&alarms, fmt.Sprintf(getOrphanAlarms, s.Context.Config.IntervalForOrphan))
	checkErr(err, "Select Orphan Alarm Failed")

	if len(alarms) > 0 {
		printOrphanAlarm(alarms)
		s.executeMoveAndRemove(alarms)

	} else {
		logger.Debug("Orphan alarm not exists.")
	}

}

// History로 이동
// Current 제거
func (s *AlarmService) executeMoveAndRemove(alarms []model.AlarmCurrent) {
	txDB, _ := s.DB.DB.BeginTransaction()
	var idStr string = makeIdsForQuery(alarms)
	logger.Debugf("ID list %s", idStr)

	_, err := txDB.GetClient().Exec(fmt.Sprintf(moveOrphanAlarms, idStr))
	checkErr(err, "move current alarm to History failed")
	if err != nil {
		txDB.Rollback()
		return
	}
	_, err = txDB.GetClient().Exec(fmt.Sprintf(removeOrphanAlarms, idStr))
	checkErr(err, "delete current alarm failed")
	if err != nil {
		txDB.Rollback()
		return
	}

	txDB.Commit()
}

func makeIdsForQuery(alarms []model.AlarmCurrent) string {
	var b bytes.Buffer
	var max int = len(alarms)
	b.WriteString("'")
	for i, alarm := range alarms {
		b.WriteString(alarm.AlarmId)
		if i <= max-1 {
			b.WriteString("','")
		}
	}
	b.WriteString("'")
	return b.String()
}

func printOrphanAlarm(alarms []model.AlarmCurrent) {
	//	func printOrphanAlarm(alarms []structure.AliveAlarm) {
	curr := time.Now().Truncate(60 * time.Second)
	timestamp := curr.UTC()
	logger.Infof("========================[ Orphan Alarms : %v ]================================", timestamp.String())
	summary := make(map[string]int)
	for _, alarm := range alarms {
		logger.Infof("  - Orphan alarm :: %s %s %v", alarm.ClusterId, alarm.AlarmId, alarm.Timestamp)
		summary[alarm.ClusterId] += 1
	}
	for key, count := range summary {
		logger.Infof("  * Cluster[%s] has %d orphan alarms", key, count)
	}
	logger.Debug("=========================================================================")
}
