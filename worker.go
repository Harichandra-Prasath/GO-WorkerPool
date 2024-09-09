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

func SpawnWorker() *Worker {
	fmt.Println("Spawning a Worker")
	return &Worker{
		Available: true,
		ID:        uuid.New(),
	}
}

func (w *Worker) Do(wg *sync.WaitGroup) {
	defer wg.Done()

	w.Available = false

	for i := 0; i <= 3; i++ {
		fmt.Printf("Worker with ID: %s Working\n", w.ID.String())
		time.Sleep(1 * time.Second)
	}

	w.Available = true
}
