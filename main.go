package main

import (
	"time"
)

func main() {
	pool := GetPool(withInitWorkers(3), withMaxWorkers(10))
	pool.Start()
	for i := 0; i < 3; i++ {
		pool.AddJob()
		time.Sleep(1 * time.Second)
	}
	pool.KillChan <- struct{}{}
	time.Sleep(20 * time.Second)
	for {
		pool.AddJob()
		time.Sleep(2 * time.Second)
	}
}
