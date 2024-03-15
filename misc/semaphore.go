package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrNoTickets = errors.New("semaWait deadline exceeded!")

// nonbinary/counting semaphores
// used in reader-writer locks where we have multiple readers and a single writer.
type Semaphore struct {
	sem     chan struct{ id int }
	timeout time.Duration // how long to wait to acquire the semaphore before giving up
}

// push to the chan on acquire denoting we are utilising the available resource
func (s *Semaphore) semaAcquire(id int) error {
	timer := time.NewTimer(s.timeout)

	select {
	case s.sem <- struct{ id int }{id: id}:
		timer.Stop()
		return nil
	case <-timer.C:
		fmt.Println("deadline exceeded in acquiring the semaphore!")
		return ErrNoTickets
	}
}

// just consume from the channel
func (s *Semaphore) semaRelease() {
	ID := <-s.sem
	fmt.Printf("releasing the lock held by :%d\n", ID.id)
}

func (s *Semaphore) IsEmpty() bool {
	return len(s.sem) == 0
}

func semaInit(count int, timeout time.Duration) *Semaphore {
	sema := &Semaphore{
		sem:     make(chan struct{ id int }, count),
		timeout: timeout,
	}
	return sema
}
