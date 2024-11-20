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
func GetAvailableWorker(P *Pool) *Worker {

	if P.Config.Scaling {
		if len(P.Workers) == 0 {
			fmt.Println("No workers Available")
			fmt.Println("Trying to Scale Up")

			// Scale Up to MaxWorkers
			if P.Current < P.Config.MaxWorkers {
				fmt.Println("Current number of workers is less than Max workers")
				w := SpawnWorker()
				go w.Start(&P.Wg, P.Workers)
				P.Current += 1
				P.Workers <- w
			} else {
				fmt.Println("Max Workers Limit reached.. Cannot Scale Up")
			}
		}
	}

	return <-P.Workers
}

// Check at regular intervals
func (P *Pool) PollStatus() {
Outer:
	for {
		select {
		case <-P.Poller.Ticker.C:
			// Case - 1A (No of workers is equal to Min workers)
			if len(P.Workers) == P.Config.MinWorkers {
				// Minimality reached
				fmt.Println("Current number of Workers is equal to MinWorkers")
				fmt.Println("Taking No action")
				continue Outer
			} else if len(P.Workers) > P.Config.MinWorkers {
				// Case - 1B
				// More workers than needed workers, Scale down the worker
				fmt.Println("More number of Workers than MinWorkers")
				worker := <-P.Workers
				fmt.Println("Removing the Worker with id:", worker.ID)
				worker.KillChan <- struct{}{}
				P.Current -= 1
				continue Outer
			}

		case <-P.Poller.quitCh:
			fmt.Println("Stopping the Poller")
			// TODO: Maybe Removing the workers from the pool??
			break Outer
		}
	}

}
