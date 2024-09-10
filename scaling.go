package main

import (
	"fmt"
	"time"
)

type Poller struct {
	Ticker *time.Ticker
	quitCh chan struct{}
}

// check for the incoming job
func GetAvailableWorkers(P *Pool) ([]*Worker, error) {
	var avail_workers []*Worker

	for _, worker := range P.Workers {
		if worker.Available {
			avail_workers = append(avail_workers, worker)
		}
	}

	if len(avail_workers) == 0 {
		fmt.Println("No workers Available")
		fmt.Println("Trying to Scale Up")

		// Scale Up to MaxWorkers
		if len(P.Workers) < P.Config.MaxWorkers {
			fmt.Println("Current number of workers is less than Max workers")
			w := SpawnWorker()
			P.Workers = append(P.Workers, w)
			avail_workers = append(avail_workers, w)
		} else {
			fmt.Println("Max Workers Limit reached.. Cannot Scale Up")
			return nil, fmt.Errorf("NO_WORKER_ERR")
		}
	}
	return avail_workers, nil
}

// Check at regular intervals
func (P *Pool) PollStatus() {

Outer:
	for {
		select {
		case <-P.Poller.Ticker.C:
			// Case - 1 (Some workers are inactive)
			for i, worker := range P.Workers {
				if worker.Available {
					// Case - 1A (No of workers is equal to Min workers)
					if len(P.Workers) == P.Config.MinWorkers {
						// Minimality reached
						fmt.Println("Current number of Workers is less than MinWorkers")
						fmt.Println("Taking No action")
						continue Outer
					} else if len(P.Workers) > P.Config.MinWorkers {
						// Case - 1B
						// More workers than needed workers, Scale down the worker
						fmt.Println("More number of Workers than MinWorkers")
						fmt.Println("Removing the Worker with id:", worker.ID)
						P.Workers = append(P.Workers[:i], P.Workers[i+1:]...)
						continue Outer
					}

				}
			}
		case <-P.Poller.quitCh:
			fmt.Println("Stopping the Poller")
			// TODO: Maybe Removing the workers from the pool??
			break Outer
		}
	}

}
