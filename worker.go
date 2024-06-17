package workerpool

import (
	"context"
	"sync"
)

type Task func(ctx context.Context) error

func NewWorker(ctx context.Context, id int64, taskChan chan Task, errorChan chan error) *Worker {
	ctx, cancel := context.WithCancel(ctx)
	return &Worker{
		ctx:       ctx,
		cancel:    cancel,
		taskChan:  taskChan,
		errorChan: errorChan,
	}
}

type Worker struct {
	ctx       context.Context
	cancel    context.CancelFunc
	taskChan  chan Task
	errorChan chan error
}

func (worker *Worker) Stop() {
	worker.cancel()
}

func (worker *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	for {
		select {
		case <-worker.ctx.Done():
			return

		case task, ok := <-worker.taskChan:
			if worker.ctx.Err() != nil {
				return
			}
			if !ok {
				return
			}

			err = task(worker.ctx)
			if err != nil {
				worker.errorChan <- err
			}
		}
	}
}
