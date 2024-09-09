package main

import (
	"time"
)

func main() {
	pool := GetPool(withInitWorkers(3), withMaxWorkers(3))
	pool.Start()
	for i := 0; i < 4; i++ {
		go pool.AddJob()
		time.Sleep(1 * time.Second)
	}
	pool.KillChan <- struct{}{}
	time.Sleep(30 * time.Second)
	for {
		pool.AddJob()
		time.Sleep(2 * time.Second)
	}
}
