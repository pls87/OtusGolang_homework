package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var wg = sync.WaitGroup{}

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var eCounter int32
	var status int32
	taskCh := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				err := task()
				if err != nil {
					atomic.AddInt32(&eCounter, 1)
				}
			}
		}()
	}

	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			if eCounter >= int32(m) {
				atomic.AddInt32(&status, 1)
				return
			}
			taskCh <- task
		}
	}()

	wg.Wait()

	if status > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
