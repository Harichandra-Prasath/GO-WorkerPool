package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type ConfigFunc func(*PoolConfig)

type PoolConfig struct {
	InitWorkers int
	MaxWorkers  int
}

type Pool struct {
	Config   PoolConfig
	JobQueue chan struct{}
	ResQueue chan int
	Workers  []*Worker
	KillChan chan struct{}
	Wg       sync.WaitGroup
	IsAlive  bool
}

func defaultConfig() PoolConfig {
	return PoolConfig{
		InitWorkers: 2,
		MaxWorkers:  10,
	}
}

func withMaxWorkers(n int) ConfigFunc {
	return func(pc *PoolConfig) {
		pc.MaxWorkers = n
	}
}

func withInitWorkers(n int) ConfigFunc {
	return func(pc *PoolConfig) {
		pc.InitWorkers = n
	}
}

func GetPool(confs ...ConfigFunc) *Pool {
	dconf := defaultConfig()

	for _, cfn := range confs {
		cfn(&dconf)
	}

	var workers []*Worker

	for i := 0; i < dconf.InitWorkers; i++ {
		workers = append(workers, SpawnWorker())
	}

	return &Pool{
		Config:   dconf,
		JobQueue: make(chan struct{}),
		KillChan: make(chan struct{}),
		ResQueue: make(chan int),
		Workers:  workers,
		IsAlive:  true,
	}
}

func (P *Pool) Start() {
	go listen(P)
}

func listen(P *Pool) {
	for {
		select {
		case <-P.JobQueue:
			Schedule(P)
		case <-P.KillChan:
			fmt.Println("Kill Signal Recieved...Waiting for the workers to finish the Job")
			P.Wg.Wait()
			fmt.Println("Pool Killed")
			P.IsAlive = false
			return
		}

	}
}

func (P *Pool) AddJob() {
	if !P.IsAlive {
		fmt.Println("Pool is not alive. Cannot add a Job")
		return
	}
	P.JobQueue <- struct{}{}
}

func Schedule(P *Pool) {

	var avail_workers []*Worker
	for _, worker := range P.Workers {
		if worker.Available {
			avail_workers = append(avail_workers, worker)
		}
	}

	if len(avail_workers) == 0 {
		fmt.Println("No workers Available. Cannot Schedule the Job")
		return
	}

	rand_index := rand.Intn(len(avail_workers))

	worker := avail_workers[rand_index]

	fmt.Printf("Worker with ID: %s is starting the work\n", worker.ID.String())
	P.Wg.Add(1)
	go worker.Do(&P.Wg)

}
