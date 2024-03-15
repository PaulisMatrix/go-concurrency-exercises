package main

import (
	"sync"
	"testing"
	"time"
)

func TestSemaDeadlineExceeded(t *testing.T) {
	// Start off 5 goroutines with 4 tickets available
	// Each just sleeping for say 1 secs. The 5th should error out with deadline exceeded
	var wg sync.WaitGroup
	n := 4
	sema := semaInit(n, 50*time.Millisecond)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			// can access the resource concurrently as long as there are n tickets available
			// for example in multiple readers, single writer, N readers can read given no writer is active
			defer wg.Done()

			// acquire the resource
			if err := sema.semaAcquire(id); err != nil {
				t.Error(err)
				return
			}

			// do the work
			time.Sleep(2 * time.Second)

			// release the semaphore
			sema.semaRelease()

		}(i + 1)
	}

	time.Sleep(1 * time.Second)
	// start off the last one a little late, ideally this should timeout since all tickets are already consumed and being worked out.

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := sema.semaAcquire(5); err != ErrNoTickets {
			t.Error(err)
			return
		}
	}()

	wg.Wait()

}
