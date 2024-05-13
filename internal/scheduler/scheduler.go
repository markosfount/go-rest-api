package scheduler

import (
	"fmt"
	"sync"
	"time"
)

func Run(interval time.Duration, done <-chan bool, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Second * interval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Running background process")
		case <-done:
			fmt.Println("Shutting down background process")
			time.Sleep(10 * time.Second)
			ticker.Stop()
			wg.Done()
			return
		}
	}
}
