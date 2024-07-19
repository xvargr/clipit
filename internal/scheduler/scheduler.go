package scheduler

import (
	"time"
)

// scheduler is responsible for running tasks at a specified interval
func Register(t time.Duration, task func()) {
	scheduler := time.NewTicker(t)

	go func() {
		defer scheduler.Stop()
		for {
			<-scheduler.C
			task()
		}
	}()
}
