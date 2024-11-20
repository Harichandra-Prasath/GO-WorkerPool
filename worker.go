package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	ID       uuid.UUID
	JobChan  chan *Job
	KillChan chan struct{}
}

type Job struct {
	ID       uuid.UUID
	WorkTime int `json:"work_time"`
}

func SpawnWorker() *Worker {

	id := uuid.New()
	fmt.Println("Worker Spawned with ID:", id)

	return &Worker{
		ID:       id,
		JobChan:  make(chan *Job),
		KillChan: make(chan struct{}),
	}
}

func (w *Worker) Start(wg *sync.WaitGroup, workers chan *Worker) {
outer:
	for {
		select {
		case j := <-w.JobChan:
			fmt.Printf("Worker with ID: %s executing the Job with ID: %s\n", w.ID, j.ID)
			Execute(j)
			fmt.Printf("Worker with ID: %s Finished the Job with ID: %s\n", w.ID, j.ID)
			wg.Done()

			workers <- w

		case <-w.KillChan:
			fmt.Printf("Woker with ID: %s Killed\n", w.ID)
			break outer
		}
	}
}

func Execute(j *Job) {

	// Mimic Work time
	time.Sleep(time.Duration(j.WorkTime) * time.Millisecond)

}
