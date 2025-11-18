# go-worker-pool

Deadly simple worker pool

## Usage

```go
package main

import (
	"context"
	"errors"
	"fmt"
	workerpool "github.com/zelenin/go-worker-pool"
	"log"
)

func main() {
	pool := workerpool.NewPool(2, 10)

	go func() {
		errorChan := pool.Errors()

		for {
			err, ok := <-errorChan
			if !ok {
				break
			}

			log.Printf("%s", err)
		}
	}()

	for i := int64(1); i < 100; i++ {
		id := i	
		pool.AddTask(func(ctx context.Context) error {
			if id%2 == 0 {
				return fmt.Errorf("Task #%d: %w", id, errors.New("task error"))
			}

			return nil
		})
	}
	
	pool.Wait()
}
```

## Author

[Aleksandr Zelenin](https://github.com/zelenin/), e-mail: [aleksandr@zelenin.me](mailto:aleksandr@zelenin.me)
