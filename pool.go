package workerpool

import (
	"sync"
)

type Pool struct {
	workersCount int64
	taskChan     chan *Task
	errorChan    chan error
	stopped      chan bool
}

func NewPool(workersCount int64, capacity int64) *Pool {
	pool := &Pool{
		workersCount: workersCount,
		taskChan:     make(chan *Task, capacity),
		errorChan:    make(chan error, capacity),
		stopped:      make(chan bool),
	}

	go pool.start()

	return pool
}

func (pool *Pool) start() {
	var wg sync.WaitGroup
	for id := int64(1); id <= pool.workersCount; id++ {
		wg.Add(1)
		worker := NewWorker(id, pool.taskChan, pool.errorChan)
		go worker.Run(&wg)
	}
	wg.Wait()
	pool.stopped <- true
}

func (pool *Pool) Wait() {
	//close(pool.taskChan)
	<-pool.stopped
	close(pool.errorChan)
	close(pool.stopped)
}

func (pool *Pool) Close() {
	close(pool.taskChan)
}

func (pool *Pool) AddTask(task *Task) {
	pool.taskChan <- task
}

func (pool *Pool) Errors() chan error {
	return pool.errorChan
}
