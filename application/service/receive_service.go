package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	alarmlog "github.com/cocktailcloud/acloud-alarm-collector/application/logger"
	"github.com/cocktailcloud/acloud-alarm-collector/application/structure"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/types"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/utils"
	"github.com/valyala/fastjson"
)

// 전달 받은 알람을 Flow(절차)에 따라 처리하고 결과를 반환합니다.
// 1. AlertFiltering
// 2. 알람 Object로 변경
// 3. 메인 Queue에 전달 합니다.
func (s *AlarmService) ReceiveAlarmsProcess(envelop *types.TransmitEnvelop) ([][]byte, error) {

	curr := time.Now().Truncate(60 * time.Second)
	clusterId := envelop.Cluster
	groupLabel, alerts, err := getAlertArray(envelop.Alerts)
	if err != nil {
		return make([][]byte, 0), nil
	}

	for _, alert := range alerts {
		labels, err := getLableValues(alert.Get("labels"))
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		// 하드코딩으로 1차 작업을 진행합니다.
		if ok, reason := checkFiltering(labels); ok {
			logger.Infof("This alarm is removed by the filter. (reason: %s)", reason)
			continue
		}

		namespace := labels["namespace"]
		workload_type := labels["workload_type"]
		workload := labels["workload"]
		pod := labels["pod"]
		container := labels["container"]
		alertname := labels["alertname"]
		severity := labels["severity"]

		description := getDescription(alert.Get("annotations"))
		//
		ruleID := getRuleId(alert.Get("annotations"))
		//
		status, err := getStringValue(alert, false, "status")
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		startAtString, err := getStringValue(alert, false, "startsAt")
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		startsAt, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", startAtString)
		if err != nil {
			logger.Errorf("invalid time: %s, %s", startAtString, err.Error())
			continue
		}
		startsAtUTC := startsAt.UTC()
		logger.Infof("recevied alarm: %s, %v, %v, %v", clusterId, alertname, status, startsAt)

		if strings.Compare(severity, "none") == 0 {
			continue
		}

		fingerprint, err := getStringValue(alert, false, "fingerprint")
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		alarmKey := fmt.Sprintf("%s-%s", fingerprint, startsAtUTC.Format(utils.DefaultTimeFormat))

		endsAtString, _ := getStringValue(alert, true, "endsAt")
		endsAtUTC := time.Time{}
		if len(endsAtString) > 0 && strings.Compare(endsAtString, "0001-01-01T00:00:00Z") != 0 {
			endsAt, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", endsAtString)
			if err == nil {
				endsAtUTC = endsAt.UTC()
			}
		}
		timestamp := curr.UTC()

		alarmObject := &structure.Alarm{
			ClusterId:    clusterId,
			GroupKey:     groupLabel,
			Fingerprint:  fingerprint,
			Namespace:    namespace,
			WorkloadType: workload_type,
			Workload:     workload,
			Pod:          pod,
			Container:    container,
			Alertname:    alertname,
			Severity:     severity,
			Description:  description,
			RuleId:       ruleID,
			Status:       status,
			AlarmKey:     alarmKey,
			StartsAt:     startsAtUTC,
			EndsAt:       endsAtUTC,
			Timestamp:    timestamp,
		}

		alertItem := structure.AlarmItem{
			Retry: 0,
			Alarm: alarmObject,
		}

		alarmhistroy := &structure.AlarmHistory{
			ClusterId:   clusterId,
			Alertname:   alertname,
			RuleId:      ruleID,
			Severity:    severity,
			Status:      status,
			StartsAt:    startsAt,
			EndsAt:      endsAtUTC,
			Description: description,
		}
		// alarmlog.History(alarmHistoryFormat(alarmObject, "json"))
		alarmlog.AlarmHistoryFormat("json", alarmhistroy)

		s.MainQueue.Add(alertItem)
	}

	return make([][]byte, 0), nil
}

func getAlertArray(alertbytes []byte) (string, []*fastjson.Value, error) {
	env, err := fastjson.ParseBytes(alertbytes)
	if err != nil {
		logger.Error("fail to parse alert: ", err.Error())
		return "", nil, err
	}
	groups := getGroupLabelFromAlarmJson(env)
	alerts := env.GetArray("alerts")
	if len(alerts) == 0 {
		logger.Warn("there is no alert")
		return "", nil, errors.New("there is no alert")
	}
	return groups, alerts, nil
}

// 알람을 전달한 그룹정보를 조회합니다.
func getGroupLabelFromAlarmJson(group *fastjson.Value) string {
	b := group.GetStringBytes("groupKey")
	if len(b) == 0 {
		return "nogroup"
	}
	return string(b)
}

func getLableValues(labels *fastjson.Value) (map[string]string, error) {
	if labels == nil {
		return nil, fmt.Errorf("'labels' filed not exists")
	}
	values := make(map[string]string)
	for i := 0; i < len(labelFields); i++ {
		r := labels.GetStringBytes(labelFields[i].Name)
		if len(r) == 0 {
			if !labelFields[i].Nullable {
				return nil, fmt.Errorf("field not exists: %s", labelFields[i].Name)
			}
			// values = append(values, sql.NullString{})
			values[labelFields[i].Name] = ""
		} else {
			// values = append(values, string(r))
			values[labelFields[i].Name] = string(r)
		}
	}

	return values, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Rule 1
//   - alarmname == 'EventWarning'
//   - Failed 문자열 포함한 경우
//   - Evicted
//   - Rebooted
//   - NodeNotReady
func checkFiltering(labels map[string]string) (bool, string) {

	allowListForEventWarning := []string{"Evicted", "Rebooted", "NodeNotReady"}
	if labels["alertname"] == "EventWarning" {
		if contains(allowListForEventWarning, labels["reason"]) {
			return false, ""
		}
		if strings.Contains(labels["reason"], "Failed") {
			return false, ""
		}
		return true, fmt.Sprintf("EventWarning(%s)", labels["reason"])
	}
	return false, ""
}
