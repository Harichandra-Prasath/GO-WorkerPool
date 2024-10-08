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
	WorkTime int `json:"work_time"`
}

func SpawnWorker() *Worker {

	id := uuid.New()
	fmt.Println("Worker Spawned with ID:", id)

	return &Worker{
		Available: true,
		ID:        id,
	}
}

func (w *Worker) Do(wg *sync.WaitGroup, j *Job) {
	defer wg.Done()

	fmt.Printf("Worker with ID: %s executing the Job with ID: %s\n", w.ID, j.ID)

	// Mimic Work time
	time.Sleep(time.Duration(j.WorkTime) * time.Second)

	fmt.Printf("Worker with ID: %s Finished the Job with ID: %s\n", w.ID, j.ID)
}
