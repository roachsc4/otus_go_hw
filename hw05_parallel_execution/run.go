package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// errorsCounter encapsulates processing of errors count.
type errorsCounter struct {
	mu         sync.Mutex
	counter    int
	errorLimit int
}

func (ec *errorsCounter) HasReachedLimit() bool {
	defer ec.mu.Unlock()
	ec.mu.Lock()
	return ec.counter > ec.errorLimit
}

func (ec *errorsCounter) Increase() {
	defer ec.mu.Unlock()
	ec.mu.Lock()
	ec.counter++
}

type Worker struct {
	wg            *sync.WaitGroup
	tasksChannel  chan Task
	errorsCounter *errorsCounter
}

// Work processes tasks in infinite loop until the given channel is closed or errors limit is reached.
func (w Worker) Work() {
	for {
		if w.errorsCounter.HasReachedLimit() {
			break
		}
		task, channelIsOpen := <-w.tasksChannel
		if !channelIsOpen {
			break
		}
		taskError := task()
		if taskError != nil {
			w.errorsCounter.Increase()
		}
	}
	w.wg.Done()
}

// putTasksIntoChannel puts all tasks one by one into the given channel and then closes it.
func putTasksIntoChannel(tasks []Task, channel chan<- Task) {
	for _, task := range tasks {
		channel <- task
	}
	close(channel)
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, N int, M int) error { //nolint:gocritic
	errorsCounter := &errorsCounter{
		counter:    0,
		errorLimit: M,
	}

	tasksChannel := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}
	wg.Add(N)
	// start N workers
	for i := 0; i < N; i++ {
		worker := Worker{
			wg:            &wg,
			tasksChannel:  tasksChannel,
			errorsCounter: errorsCounter,
		}
		go worker.Work()
	}

	// put all tasks to the channel and wait for the end of processing
	go putTasksIntoChannel(tasks, tasksChannel)
	wg.Wait()
	if errorsCounter.HasReachedLimit() {
		return ErrErrorsLimitExceeded
	}
	return nil
}
