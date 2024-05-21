package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	interval time.Duration
	done     <-chan bool
	wg       *sync.WaitGroup
}

func NewScheduler(interval time.Duration, done <-chan bool, wg *sync.WaitGroup) *Scheduler {
	return &Scheduler{
		interval: interval,
		done:     done,
		wg:       wg,
	}
}

func (s *Scheduler) Run() {
	ticker := time.NewTicker(time.Second * s.interval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Running background process")
		case <-s.done:
			fmt.Println("Shutting down background process")
			time.Sleep(10 * time.Second)
			ticker.Stop()
			s.wg.Done()
			return
		}
	}
}
