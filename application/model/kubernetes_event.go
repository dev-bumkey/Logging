package model

import "time"

type KubernetesEvent struct {
	ClusterId    string    `json:"clusterid" db:"cluster_id"`
	EventId      string    `json:"event_id" db:"event_id"`
	EventType    string    `json:"event_type" db:"event_type"`
	Node         string    `json:"node" db:"node"`
	Namespace    string    `json:"namespace" db:"namespace"`
	Workload     string    `json:"workload" db:"workload"`
	Pod          string    `json:"pod" db:"pod"`
	Container    string    `json:"container" db:"container"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
	FirstTime    time.Time `json:"first_time" db:"first_time"`
	LastTime     time.Time `json:"last_time" db:"last_time"`
	Reason       string    `json:"reason" db:"reason"`
	Message      string    `json:"message" db:"message"`
	Count        string    `json:"count" db:"count"`
	WorkloadType string    `json:"workload_type" db:"workload_type"`
}
