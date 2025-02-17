package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task)
	resultsCh := make(chan error, len(tasks))
	closeCh := make(chan struct{})
	wg := new(sync.WaitGroup)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task := <-tasksCh:
					resultsCh <- task()
				case <-closeCh:
					return
				}
			}
		}()
	}

	for _, task := range tasks {
		tasksCh <- task
	}

	err := func() error {
		var errCnt, allCnt int
		defer close(closeCh)

		for {
			allCnt++
			err := <-resultsCh
			if err != nil {
				errCnt++
			}

			if allCnt == len(tasks) {
				return nil
			} else if errCnt >= m {
				return ErrErrorsLimitExceeded
			}
		}
	}()

	wg.Wait()
	return err
}
