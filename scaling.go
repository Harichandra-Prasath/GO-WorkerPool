package main

import "fmt"

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
