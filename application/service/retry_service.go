package service

import "github.com/cocktailcloud/acloud-monitoring-common/v2/logger"

const (
	retryLimitMessage = `
	This alarm is cleared because the number of retries has been exceeded.
		- cluster_id : %s
		- alarm_name : %s
		- namespace : %s
		- workload : %s
		- severity : %s
		- starts_at : %v
		- description : %s
	`
)

// 재시도 큐에 존재하는 알람을 처리 합니다. (실패시 시간 term을 주기 위해 1분마다 재시도 합니다.)
func (s *AlarmService) RetryProcess() {

	if s.RetryQueue.IsEmpty() {
		// logger.Debug("Retry Queue is empty")
		return
	}
	contents, ok := s.RetryQueue.Poll()
	logger.Debugf("retry entry is %d'th: %v", contents.Retry, contents.Alarm.Fingerprint)
	if ok {
		if contents.Retry > s.Context.Config.RetryLimit {
			// 이 알람은 재시도 횟수를 초과 하였기 때문에 제거 됩니다.
			logger.Errorf(retryLimitMessage,
				contents.Alarm.ClusterId,
				contents.Alarm.Alertname,
				contents.Alarm.Namespace,
				contents.Alarm.Workload,
				contents.Alarm.Severity,
				contents.Alarm.StartsAt,
				contents.Alarm.Description,
			)
			return
		}
		err := s.process(contents)
		if err != nil {
			logger.Errorf("fail to retry %v", err)
			contents.Retry += 1
			s.RetryQueue.Add(contents)
		}

	}
}
