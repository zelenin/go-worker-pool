package workerpool

import (
	"fmt"
	"sync"
)

type TaskError struct {
	Id    int64
	error error
}

func (err TaskError) Error() string {
	return fmt.Sprintf("Task #%d: %s", err.Id, err.error)
}

func (err TaskError) Unwrap() error {
	return err.error
}

type Processer func(id int64) error

type Task struct {
	Id        int64
	Processer Processer
}

func NewTask(id int64, processer Processer) *Task {
	return &Task{
		Id:        id,
		Processer: processer,
	}
}

func NewWorker(id int64, taskChan chan *Task, errorChan chan error) *Worker {
	return &Worker{
		id:        id,
		taskChan:  taskChan,
		errorChan: errorChan,
	}
}

type Worker struct {
	id        int64
	taskChan  chan *Task
	errorChan chan error
}

func (worker *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	for task := range worker.taskChan {
		err = task.Processer(task.Id)
		if err != nil {
			worker.errorChan <- TaskError{
				Id:    task.Id,
				error: err,
			}
		}
	}
}
