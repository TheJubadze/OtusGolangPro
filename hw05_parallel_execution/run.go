package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var mu sync.Mutex
	var wg sync.WaitGroup
	mc := Counter{m, &mu}
	tasksChan := make(chan Task)
	shouldCountErrors := m > 0
	go enqueueTasks(&tasks, tasksChan, &mc, shouldCountErrors)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go run(tasksChan, &wg, &mc)
	}
	wg.Wait()
	if shouldCountErrors && mc.Get() <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func enqueueTasks(tasks *[]Task, tasksChan chan Task, m *Counter, shouldCountErrors bool) {
	defer close(tasksChan)
	for _, task := range *tasks {
		if shouldCountErrors && m.Get() <= 0 {
			return
		}
		tasksChan <- task
	}
}

func run(ch chan Task, wg *sync.WaitGroup, m *Counter) {
	defer wg.Done()
	for task := range ch {
		if task() != nil {
			m.Dec()
		}
	}
}

type Counter struct {
	value int
	mu    *sync.Mutex
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Dec() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
}

func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
