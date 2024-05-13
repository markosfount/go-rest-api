package scheduler

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func Run(interval time.Duration, quit <-chan os.Signal, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Second * interval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Running background process")
		case _, ok := <-quit:
			if !ok {
				fmt.Println("Shutting down background process")
				time.Sleep(10 * time.Second)
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}
}
