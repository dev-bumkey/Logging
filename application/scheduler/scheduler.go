package scheduler

type Scheduler interface {
	Run() error
	Stop() error
}
