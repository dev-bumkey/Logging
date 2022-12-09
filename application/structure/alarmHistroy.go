package structure

import (
	"time"
)

type AlarmHistory struct {
	ClusterId   string    `json:"cluster_Id"`
	Alertname   string    `json:"alter_Name"`
	RuleId      string    `json:"rule_ID"`
	Severity    string    `json:"serverity"`
	Status      string    `json:"status"`
	StartsAt    time.Time `json:"start_At"`
	EndsAt      time.Time `json:"end_At"`
	Description string    `json:"decription"`
}
