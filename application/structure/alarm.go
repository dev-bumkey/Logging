package structure

import (
	"time"
)

type AlarmItem struct {
	Retry int
	Alarm *Alarm
}

type Alarm struct {
	ClusterId    string
	GroupKey     string
	Fingerprint  string
	Namespace    string
	WorkloadType string
	Workload     string
	Pod          string
	Container    string
	Alertname    string
	Severity     string
	Description  string
	//
	RuleId string
	//
	Status    string
	AlarmKey  string
	StartsAt  time.Time
	EndsAt    time.Time
	Timestamp time.Time
}
