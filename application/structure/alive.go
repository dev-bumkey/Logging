package structure

import "time"

type AliveAlarm struct {
	ClusterId string    `db:"cluster_id"`
	AlarmId   string    `db:"alarm_id"`
	Status    string    `db:"status"`
	Timestamp time.Time `db:"timestamp"`
}
