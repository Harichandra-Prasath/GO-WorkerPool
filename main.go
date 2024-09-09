package main

import (
	"time"
)

func main() {

	confs := []ConfigFunc{withInitWorkers(3), withMinWorkers(2), withPollPeriod(5), withMaxWorkers(5)}

	pool := GetPool(confs...)
	pool.Start()

	go pool.PollStatus()

	for i := 0; i < 5; i++ {
		go pool.AddJob()
	}
	time.Sleep(30 * time.Second)
	for {
		pool.AddJob()
		time.Sleep(2 * time.Second)
	}
}
