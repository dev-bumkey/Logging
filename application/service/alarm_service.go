package service

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cocktailcloud/acloud-alarm-collector/application/model"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/types"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/utils"
	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

var labelFields = []struct {
	Name     string
	Nullable bool
	Value    interface{}
}{
	{"namespace", true, nil},
	{"workload_type", true, nil},
	{"workload", true, nil},
	{"pod", true, nil},
	{"container", true, nil},
	{"alertname", false, nil},
	{"severity", true, nil},
	{"reason", true, nil},
}

const (
	deleteAlarmCurrent = "delete from alarm_current where cluster_id=$1 and group_label=$2"
	countAlarmCurrent  = "select count(1) from alarm_current where cluster_id=$1 and group_label=$2"
)

// 기존 작업 프로세스
// 1. 현재 활성화 알람 전체 제거
// 2. 알람을 이력 테이블에 입력
// 3. 알람을 활성화 테이블에 입력
func (service *AlarmService) InsertAlerts(envelop *types.TransmitEnvelop) ([][]byte, error) {
	insertedAlarms := make([][]byte, 0)
	clusterId := envelop.Cluster

	groupLabel, alerts, err := getAlertArray(envelop.Alerts)
	if err != nil {
		return insertedAlarms, nil
	}

	currentBulk := make([]interface{}, 0)
	historyBulk := make([]interface{}, 0)
	curr := time.Now().Truncate(60 * time.Second)

	for _, alert := range alerts {
		alarmId := uuid.New().String()
		// namespace := getNamespaceFromLabel(alert.Get("labels"))
		labels, err := getLableValues(alert.Get("labels"))
		if err != nil {
			logger.Error(err.Error())
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

		currentModel := &model.AlarmCurrent{
			ClusterId:    clusterId,
			AlarmId:      alarmId,
			GroupLabel:   groupLabel,
			Namespace:    namespace,
			WorkloadType: workload_type,
			WorkloadName: workload,
			Pod:          pod,
			Container:    container,
			AlarmName:    alertname,
			Severity:     severity,
			Description:  description,
			Status:       status,
			Startsat:     startsAtUTC,
			AlarmKey:     alarmKey,
			Endsat:       endsAtUTC,
			Timestamp:    timestamp,
		}
		historyModel := &model.AlarmHistory{
			ClusterId:    clusterId,
			AlarmId:      alarmId,
			Namespace:    namespace,
			WorkloadType: workload_type,
			WorkloadName: workload,
			Pod:          pod,
			Container:    container,
			AlarmName:    alertname,
			Severity:     severity,
			Description:  description,
			Status:       status,
			Startsat:     startsAtUTC,
			AlarmKey:     alarmKey,
			Endsat:       endsAtUTC,
			Timestamp:    timestamp,
		}

		currentBulk = append(currentBulk, currentModel)
		historyBulk = append(historyBulk, historyModel)

		alert.Set("fingerprint", fastjson.MustParse(fmt.Sprintf(`"%s"`, alarmKey)))
		alert.Del("generatorURL")
		alert.Get("labels").Set("alarmId", fastjson.MustParse(fmt.Sprintf(`"%s"`, alarmId)))
		alert.Set("timestamp", fastjson.MustParse(fmt.Sprintf(`"%s"`, curr.UTC().Format(utils.DefaultTimeFormat))))

		insertedAlarms = append(insertedAlarms, []byte(alert.String()))
	}

	txDB, _ := service.DB.DB.BeginTransaction()

	// Delete Current
	// delCnt, err := txDB.GetClient().Delete(&model.AlarmCurrent{ClusterId: clusterId, GroupLabel: groupLabel})
	delCnt, err := txDB.GetClient().SelectInt(countAlarmCurrent, clusterId, groupLabel)
	if err != nil {
		logger.Error("Previous current alarm count error: ", err)
		txDB.Rollback()
		return make([][]byte, 0), nil
	}
	_, err = txDB.GetClient().Exec(deleteAlarmCurrent, clusterId, groupLabel)
	if err != nil {
		logger.Error("Previous current alarm delete error: ", err)
		txDB.Rollback()
		return make([][]byte, 0), nil
	}
	logger.Infof("Previous current alarm Deleted (%d) by group %s", delCnt, groupLabel)

	if len(historyBulk) > 0 {
		// Insert History
		err = txDB.GetClient().Insert(historyBulk...)
		if err != nil {
			logger.Error("history alarm insert error: ", err)
			txDB.Rollback()
			return make([][]byte, 0), nil
		}
	}

	if len(currentBulk) > 0 {
		// Insert Current
		err = txDB.GetClient().Insert(currentBulk...)
		if err != nil {
			logger.Error("current alarm insert error: ", err)
			txDB.Rollback()
			return make([][]byte, 0), nil
		}
	}

	txDB.Commit()

	return insertedAlarms, nil
}

func getRuleId(item *fastjson.Value) string {
	if item == nil {
		return ""
	}

	o, err := item.Object()
	if err != nil {
		logger.Warn(err.Error())
		return ""
	}

	messages := bytes.Buffer{}
	o.Visit(func(k []byte, v *fastjson.Value) {
		key := string(k)

		if !strings.HasPrefix(string(key), "rule") {
			return
		}
		var result = strings.Trim(key, "rule_id")
		messages.WriteString(result)
		messages.Write(v.GetStringBytes())

	})

	return messages.String()
}

func getEngMsg(item *fastjson.Value) string {
	if item == nil {
		return ""
	}

	o, err := item.Object()
	if err != nil {
		logger.Warn(err.Error())
		return ""
	}

	messages := bytes.Buffer{}
	o.Visit(func(k []byte, v *fastjson.Value) {
		key := string(k)

		if !strings.HasPrefix(string(key), "message_en") {
			return
		}
		var result = strings.Trim(key, "message_en")
		messages.WriteString(result)
		messages.Write(v.GetStringBytes())

	})

	return messages.String()
}

func getDescription(item *fastjson.Value) string {
	if item == nil {
		return ""
	}

	o, err := item.Object()
	if err != nil {
		logger.Warn(err.Error())
		return ""
	}

	messages := bytes.Buffer{}
	o.Visit(func(k []byte, v *fastjson.Value) {
		key := string(k)
		if !strings.HasPrefix(string(key), "message") {
			return
		}

		if messages.Len() == 0 {
			messages.WriteString("{")
		} else {
			messages.WriteString(",")
		}
		messages.WriteString("\"")
		messages.WriteString(key)
		messages.WriteString("\": ")
		messages.WriteString("\"")
		messages.Write(v.GetStringBytes())
		messages.WriteString("\"")
	})

	if messages.Len() > 0 {
		messages.WriteString("}")
	}

	return messages.String()
}

func getStringValue(item *fastjson.Value, nullable bool, keys ...string) (string, error) {
	byteValue := item.GetStringBytes(keys...)
	if byteValue == nil {
		if nullable {
			return "", nil
		} else {
			return "", fmt.Errorf("field not exists: %s", keys)
		}
	} else {
		return string(byteValue), nil
	}
}
