package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	workersCount int64
	workers      []*Worker
	taskChan     chan Task
	errorChan    chan error
	stoppedChan  chan bool
	mu           sync.Mutex
	stopped      bool
}

func NewPool(workersCount int64, capacity int64) *Pool {
	pool := &Pool{
		workersCount: workersCount,
		workers:      make([]*Worker, 0, workersCount),
		taskChan:     make(chan Task, capacity),
		errorChan:    make(chan error, capacity),
		stoppedChan:  make(chan bool),
		mu:           sync.Mutex{},
		stopped:      false,
	}

	go pool.start()

	return pool
}

func (pool *Pool) start() {
	var wg sync.WaitGroup
	for id := int64(1); id <= pool.workersCount; id++ {
		wg.Add(1)
		ctx := context.Background()
		worker := NewWorker(ctx, id, pool.taskChan, pool.errorChan)
		pool.workers = append(pool.workers, worker)
		go worker.Run(&wg)
	}
	wg.Wait()
	pool.stoppedChan <- true
}

func (pool *Pool) Wait() {
	close(pool.taskChan)
	<-pool.stoppedChan
	close(pool.errorChan)
	close(pool.stoppedChan)
}

func (pool *Pool) Stopped() bool {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.stopped
}

func (pool *Pool) Stop() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for _, worker := range pool.workers {
		worker.Stop()
	}
	pool.stopped = true
}

func (pool *Pool) AddTask(task Task) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.stopped {
		return
	}
	pool.taskChan <- task
}

func (pool *Pool) Errors() chan error {
	return pool.errorChan
}
