package utils

import (
	"github.com/cocktailcloud/acloud-alarm-collector/application/structure"
	"github.com/eapache/queue"
)

type AlarmQueue struct {
	Queue *queue.Queue
}

func NewQueue() AlarmQueue {
	return AlarmQueue{
		Queue: queue.New(),
	}
}

func (q *AlarmQueue) Add(elem structure.AlarmItem) {
	q.Queue.Add(elem)
}

func (q *AlarmQueue) Poll() (structure.AlarmItem, bool) {
	if q.Queue.Length() > 0 {
		return q.Queue.Remove().(structure.AlarmItem), true
	}
	return structure.AlarmItem{}, false
}

func (q *AlarmQueue) Clear() {
	for q.Queue.Length() > 0 {
		q.Queue.Remove()
	}
}

func (q *AlarmQueue) IsEmpty() bool {
	return q.Queue.Length() == 0
}
