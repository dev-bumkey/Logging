package model

import "time"

type AlarmHistory struct {
	ClusterId     string    `json:"clusterid" db:"cluster_id"`
	AlarmId       string    `json:"alarm_id" db:"alarm_id"`
	Namespace     string    `json:"namespace" db:"namespace"`
	WorkloadType  string    `json:"workload_type" db:"workload_type"`
	WorkloadName  string    `json:"workload_name" db:"workload_name"`
	Pod           string    `json:"pod" db:"pod"`
	Container     string    `json:"container" db:"container"`
	AlarmKey      string    `json:"alarm_key" db:"alarm_key"`
	AlarmName     string    `json:"alarm_name" db:"alarm_name"`
	Severity      string    `json:"severity" db:"severity"`
	Description   string    `json:"description" db:"description"`
	Status        string    `json:"status" db:"status"`
	Startsat      time.Time `json:"startsat" db:"startsat"`
	Endsat        time.Time `json:"endsat" db:"endsat"`
	ProcessStatus int       `json:"process_status" db:"process_status"`
	ErrorDetail   string    `json:"error_detail" db:"error_detail"`
	SendsAt       time.Time `json:"sends_at" db:"sends_at"`
	AckedAt       time.Time `json:"acked_at" db:"acked_at"`
	AckChunk      string    `json:"ack_chunk" db:"ack_chunk"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
}
