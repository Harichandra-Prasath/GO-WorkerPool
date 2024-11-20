package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

func Bench() {

	// Starting M:M Jobs to Goroutines

	WORK_TIME := 500
	N_JOBS := 1000000

	start_m := time.Now()

	wg := sync.WaitGroup{}

	for i := 0; i < N_JOBS; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			defer Execute(&Job{
				ID:       uuid.New(),
				WorkTime: WORK_TIME,
			})
		}(&wg)

	}

	wg.Wait()

	fmt.Printf("Time taken to complete %d Jobs with %d go-routines:%s\n", N_JOBS, N_JOBS, time.Since(start_m).String())

	N_WORKERS := N_JOBS / 2
	confs := []ConfigFunc{withInitWorkers(N_WORKERS), withMinWorkers(0), withPollPeriod(5), withMaxWorkers(N_JOBS)}

	pool := GetPool(confs...)

	pool.Start()

	start_m = time.Now()

	for i := 0; i < N_JOBS; i++ {
		pool.AddJob(&Job{
			ID:       uuid.New(),
			WorkTime: WORK_TIME,
		})
	}

	pool.Wg.Wait()

	fmt.Printf("Time taken to complete %d Jobs with %d Workers:%s\n", N_JOBS, N_WORKERS, time.Since(start_m).String())

}
