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
	MinWorkers  int
	Poll_Period int
}

type Pool struct {
	Config   PoolConfig
	JobQueue chan struct{}
	Workers  []*Worker
	KillChan chan struct{}
	Wg       sync.WaitGroup
	IsAlive  bool
}

func defaultConfig() PoolConfig {
	return PoolConfig{
		InitWorkers: 2,
		MaxWorkers:  4,
		MinWorkers:  1,
		Poll_Period: 10,
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

func withPollPeriod(n int) ConfigFunc {
	return func(pc *PoolConfig) {
		pc.Poll_Period = n
	}
}

func withMinWorkers(n int) ConfigFunc {
	return func(pc *PoolConfig) {
		pc.MinWorkers = n
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

	avail_workers, err := GetAvailableWorkers(P)
	if err != nil {
		fmt.Println("Error in getting available Workers:", err)
		return
	}
	rand_index := rand.Intn(len(avail_workers))

	worker := avail_workers[rand_index]

	P.Wg.Add(1)
	go worker.Do(&P.Wg)

}
