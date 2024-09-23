package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	Available bool
	ID        uuid.UUID
}

type Job struct {
	ID       uuid.UUID
	WorkTime int
}

func SpawnWorker() *Worker {
	fmt.Println("Spawning a Worker")
	return &Worker{
		Available: true,
		ID:        uuid.New(),
	}
}

func (w *Worker) Do(wg *sync.WaitGroup, j *Job) {
	defer wg.Done()

	fmt.Printf("Worker with ID: %s executing the Job with ID: %s\n", w.ID, j.ID)

	// Mimic Work time
	time.Sleep(time.Duration(j.WorkTime) * time.Second)

	fmt.Printf("Worker with ID: %s Finished the Job with ID: %s\n", w.ID, j.ID)
}
