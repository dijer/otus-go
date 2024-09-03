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
	errCount := int64(0)
	wg := sync.WaitGroup{}
	ch := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for task := range ch {
				err := task()
				if err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if int(atomic.LoadInt64(&errCount)) >= m && m > 0 {
			break
		}
		ch <- task
	}
	close(ch)

	wg.Wait()

	if int(atomic.LoadInt64(&errCount)) >= m && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
