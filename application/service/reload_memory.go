package service

import (
	"fmt"

	"github.com/cocktailcloud/acloud-alarm-collector/application/structure"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
)

// DB로 부터 활성화 알람을 조회하여 메모리 데이터를 replace 합니다.
func (s *AlarmService) ReloadMemory() {
	var alarms []structure.AliveAlarm
	exec := s.DB.DB.GetClient()
	_, err := exec.Select(&alarms, "select cluster_id, alarm_id, status, timestamp from alarm_current where status='firing'")
	checkErr(err, "Select Alive Alarm Failed")
	dataByCluster := make(map[string](map[string]structure.AliveAlarm))
	for _, alarm := range alarms {
		arr, ok := dataByCluster[alarm.ClusterId]
		if !ok {
			arr = make(map[string]structure.AliveAlarm)
			dataByCluster[alarm.ClusterId] = arr
		}
		arr[alarm.AlarmId] = alarm
	}
	copiedList := s.Context.GetAliveAlarm()
	printMemoryAlarm(copiedList, dataByCluster)
	s.Context.SetAliveAlarm(dataByCluster)
}

func checkErr(err error, msg string) {
	if err != nil {
		logger.Error("msg", msg, err)
	}
}

func printMemoryAlarm(prevAlarms map[string](map[string]structure.AliveAlarm), alarms map[string](map[string]structure.AliveAlarm)) {
	prevCount := make(map[string]int)
	prevItems := make(map[string]map[string]structure.AliveAlarm)
	for a, b := range prevAlarms {
		prevCount[a] = len(b)
		prevItems[a] = b
	}
	logger.Debug("========================[ Active Alarms ]================================")
	for key, items := range alarms {
		diffStr := "SAME AS ABOVE"
		if prevCount[key] != len(items) {
			diffStr = fmt.Sprintf("DIFF (prev: %d)", prevCount[key])
		}
		logger.Debugf("Cluster: %s (count: %d)-%s", key, len(items), diffStr)
		// for id, alarm := range items {
		// 	logger.Debugf("    id: %s, timestamp: %s", id, alarm.Timestamp)
		// }
	}
	logger.Debug("=========================================================================")
}

func (s *AlarmService) UpdateMemoryNew(clusterId string, alarmId string) {
	copiedList := s.Context.GetAliveAlarm()
	arr, ok := copiedList[clusterId]
	if !ok {
		arr = make(map[string]structure.AliveAlarm)
		copiedList[clusterId] = arr
	}
	arr[alarmId] = structure.AliveAlarm{
		ClusterId: clusterId,
		AlarmId:   alarmId,
	}
	s.Context.SetAliveAlarm(copiedList)
}

func (s *AlarmService) UpdateMemoryClear(clusterId string, alarmId string) {
	copiedList := s.Context.GetAliveAlarm()
	arr, ok := copiedList[clusterId]
	if ok {
		delete(arr, alarmId)
	}
	s.Context.SetAliveAlarm(copiedList)
}
