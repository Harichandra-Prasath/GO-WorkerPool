package main

import "time"

func main() {
	pool := GetPool(withInitWorkers(5), withMaxWorkers(10))
	pool.Start()
	for i := 0; i <= 10; i++ {
		pool.AddJob()
		time.Sleep(1 * time.Second)
	}
}
