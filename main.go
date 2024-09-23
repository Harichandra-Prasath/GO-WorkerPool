package main

import (
	"time"

	"github.com/google/uuid"
)

func main() {

	confs := []ConfigFunc{withInitWorkers(3), withMinWorkers(2), withPollPeriod(5), withMaxWorkers(5)}

	pool := GetPool(confs...)
	pool.Start()

	go pool.PollStatus()

	for i := 0; i < 5; i++ {
		pool.AddJob(&Job{
			ID:       uuid.New(),
			WorkTime: 30,
		})
	}

	time.Sleep(100 * time.Second)
}
