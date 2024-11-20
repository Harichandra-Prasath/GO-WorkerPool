package main

import (
	"fmt"
	"sync"
	"time"
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
	JobQueue chan *Job
	Workers  chan *Worker
	KillChan chan struct{}
	Wg       sync.WaitGroup
	Poller   *Poller
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

	workers := make(chan *Worker, dconf.InitWorkers)

	return &Pool{
		Config:   dconf,
		JobQueue: make(chan *Job),
		KillChan: make(chan struct{}),
		Workers:  workers,
		Poller: &Poller{
			Ticker: time.NewTicker(time.Duration(dconf.Poll_Period) * time.Second),
			quitCh: make(chan struct{}),
		},
		IsAlive: true,
	}
}

func (P *Pool) Start() {

	fmt.Println("Spawning the Initial Workers")

	for i := 0; i < P.Config.InitWorkers; i++ {
		P.Workers <- SpawnWorker()
	}

	go listen(P)
	go P.PollStatus()
}

func listen(P *Pool) {
	for {
		select {
		case job := <-P.JobQueue:
			Schedule(P, job)
		case <-P.KillChan:
			fmt.Println("Kill Signal Recieved...Waiting for the workers to finish the Job")
			P.Poller.quitCh <- struct{}{}
			P.Wg.Wait()
			fmt.Println("Pool Killed")
			P.IsAlive = false
			return
		}

	}
}

func (P *Pool) Kill() {
	P.KillChan <- struct{}{}
}

func (P *Pool) AddJob(J *Job) {
	if !P.IsAlive {
		fmt.Println("Pool is not alive. Cannot add a Job")
		return
	}
	P.JobQueue <- J
}

func Schedule(P *Pool, j *Job) {

	worker := GetAvailableWorker(P)

	fmt.Println("Choosen Worker:", worker.ID, " for the Job", j.ID)

	P.Wg.Add(1)
	go worker.Do(&P.Wg, j, P.Workers)

}
