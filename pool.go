package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type ConfigFunc func(*PoolConfig)

type PoolConfig struct {
	InitWorkers int
	MaxWorkers  int
}

type Pool struct {
	Config         PoolConfig
	CurrentWorkers int
	JobQueue       chan struct{}
	ResQueue       chan int
	KillChan       chan os.Signal
	Wg             sync.WaitGroup
	IsAlive        bool
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

	return &Pool{
		Config:         dconf,
		CurrentWorkers: 0,
		JobQueue:       make(chan struct{}),
		ResQueue:       make(chan int),
		KillChan:       make(chan os.Signal, 1),
		IsAlive:        true,
	}
}

func (P *Pool) Start() {

	signal.Notify(P.KillChan, os.Interrupt)
	go P.listen()
}

func (P *Pool) listen() {
	for {
		select {
		case <-P.JobQueue:
			fmt.Println("Job Recieved")
		case <-P.KillChan:
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
