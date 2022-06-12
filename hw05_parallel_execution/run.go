package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrIncorrectWorkersCount = errors.New("workers count should be more than 0")
)

type Task func() error

func validateParameters(workersCount, maxErrorsCount int) error {
	if workersCount <= 0 {
		return ErrIncorrectWorkersCount
	}
	if maxErrorsCount <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	if err := validateParameters(workersCount, maxErrorsCount); err != nil {
		return err
	}

	tasksCount := len(tasks)

	if tasksCount == 0 {
		return nil
	}

	if workersCount > tasksCount {
		workersCount = tasksCount
	}

	tasksChannel := make(chan Task, len(tasks))

	wg := &sync.WaitGroup{}

	wg.Add(workersCount)

	var errorsCount int32

	mutex := &sync.Mutex{}

	for i := 1; i <= workersCount; i++ {
		go worker(tasksChannel, wg, &errorsCount, maxErrorsCount, mutex)
	}

	for _, task := range tasks {
		tasksChannel <- task
	}

	close(tasksChannel)

	wg.Wait()

	if errorsCount >= int32(maxErrorsCount) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(tasksChan <-chan Task, wg *sync.WaitGroup, errorsCount *int32, maxErrorsCount int, m *sync.Mutex) {
	defer wg.Done()

	for {
		if task, ok := <-tasksChan; ok {
			err := task()
			m.Lock()

			if err != nil {
				*errorsCount++
			}

			if *errorsCount >= int32(maxErrorsCount) {
				m.Unlock()
				return
			}

			m.Unlock()
		} else {
			return
		}
	}
}
