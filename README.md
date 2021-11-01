# go-worker-pool

Deadly simple worker pool

## Usage

```go
package main

import (
	"errors"
	workerpool "github.com/zelenin/go-worker-pool"
	"log"
	"time"
)

func main() {
	pool := workerpool.NewPool(2, 2)

	go func() {
		log.Printf("error handler start")
		errorChan := pool.Errors()

		for {
			err, ok := <-errorChan
			if !ok {
				break
			}

			taskId := err.(workerpool.TaskError).Id
			err = errors.Unwrap(err)
			log.Printf("task #%d err: %s", taskId, err)
		}

		log.Printf("err handler finished")
	}()

	for i := int64(1); i < 100; i++ {
		log.Printf("Adding Task #%d", i)
		pool.AddTask(workerpool.NewTask(i, func(id int64) error {
			log.Printf("Task #%d started", id)
			time.Sleep(10 * time.Second)
			log.Printf("Task #%d finished", id)

			if id%2 == 0 {
				return errors.New("task error")
			}

			return nil
		}))
		log.Printf("Added Task #%d", i)
	}
	pool.Wait()
}
```

## Author

[Aleksandr Zelenin](https://github.com/zelenin/), e-mail: [aleksandr@zelenin.me](mailto:aleksandr@zelenin.me)
