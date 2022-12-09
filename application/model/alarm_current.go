package model

import "time"

type AlarmCurrent struct {
	ClusterId    string `json:"clusterid" db:"cluster_id"`
	AlarmId      string `json:"alarm_id" db:"alarm_id"`
	GroupLabel   string `json:"group_label" db:"group_label"`
	Namespace    string `json:"namespace" db:"namespace"`
	WorkloadType string `json:"workload_type" db:"workload_type"`
	WorkloadName string `json:"workload_name" db:"workload_name"`
	Pod          string `json:"pod" db:"pod"`
	Container    string `json:"container" db:"container"`
	AlarmKey     string `json:"alarm_key" db:"alarm_key"`
	AlarmName    string `json:"alarm_name" db:"alarm_name"`
	Severity     string `json:"severity" db:"severity"`
	Description  string `json:"description" db:"description"`
	//
	RuleId string `json:"ruld_Id" db:"rule_Id"`
	//
	Status    string    `json:"status" db:"status"`
	Startsat  time.Time `json:"startsat" db:"startsat"`
	Endsat    time.Time `json:"endsat" db:"endsat"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
