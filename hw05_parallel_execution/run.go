package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var eCounter, status int32
	taskCh := make(chan Task)
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if task() != nil {
					atomic.CompareAndSwapInt32(&status, atomic.AddInt32(&eCounter, 1)-int32(m), 1)
				}
			}
		}()
	}

	go func(tasks []Task) {
		defer close(taskCh)
		for _, task := range tasks {
			if atomic.CompareAndSwapInt32(&status, 1, 1) {
				return
			}
			taskCh <- task
		}
	}(tasks)

	wg.Wait()

	if atomic.CompareAndSwapInt32(&status, 1, 1) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
