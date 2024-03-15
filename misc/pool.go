package main

import (
	"fmt"
	"sync"
	"time"
)

type Job func()

type Pool struct {
	workQueue       chan Job
	deadLetterQueue chan Job
	waitGroup       sync.WaitGroup
}

func NewPool(numWorkers, numDLWorkers int) *Pool {
	p := &Pool{
		workQueue: make(chan Job),
	}

	// workers for main queue
	p.waitGroup.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		// each worker listening on the workQueue
		go func() {
			for job := range p.workQueue {
				job()
			}
			p.waitGroup.Done()
		}()
	}

	// workers for the dead letter queue
	p.waitGroup.Add(numDLWorkers)

	return p
}

func (p *Pool) AddJob(job Job) {
	p.workQueue <- job
}

func (p *Pool) PoolFin() {
	close(p.workQueue)
	p.waitGroup.Wait()
}

func WorkerPools() {
	// dummy worker pool
	numWorkers := 3
	numDLWorkers := 1

	p := NewPool(numWorkers, numDLWorkers)

	// push jobs to the queue
	for i := 0; i <= 30; i++ {
		var job Job
		job = func() {
			// simulate work
			time.Sleep(1 * time.Second)
			fmt.Println("job completed!")
		}
		p.AddJob(job)
	}

	p.PoolFin()
}
