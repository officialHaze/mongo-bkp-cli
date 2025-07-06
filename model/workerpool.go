package model

import (
	"log"
	"sync"
)

type Job func() error

type WorkerPool struct {
	TotalWorkers int
	Jobs         chan Job
	Wg           sync.WaitGroup
}

func NewWorkerPool(totalworkers int) *WorkerPool {
	return &WorkerPool{
		TotalWorkers: totalworkers,
		Jobs:         make(chan Job),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.TotalWorkers; i++ {
		wp.Wg.Add(1)
		go func() {
			// log.Printf("Starting worker %d", workerId)
			defer wp.Wg.Done()
			job, exists := <-wp.Jobs
			if !exists {
				return
			}

			err := job()
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

// Submit a job
func (wp *WorkerPool) Submit(job Job) {
	wp.Jobs <- job
}

func (wp *WorkerPool) Stop() {
	close(wp.Jobs)
	wp.Wg.Wait()
}
