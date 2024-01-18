package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// problem with this CAS operation is, the philosopher goroutines which wont be able to do a successful CAS operation will just wait burning CPU cycles?
// so don't waste CPU cycles and put them to sleep using kernel FUTEX calls

type ForkFUTEX struct {
	// set, unset this flag for CAS operationn
	// set = 1 = Availble for pickup
	// unset = 0 = Not Available for pickup
	// wait = 2 = Put to sleep by futex kernel call when flag is NA
	pickFlag uint32
}

type PhilosopherFUTEX struct {
	id    int
	left  *ForkFUTEX
	right *ForkFUTEX
}

func (p PhilosopherFUTEX) wantsToEatFUTEX(wg *sync.WaitGroup) {
	defer wg.Done()

	// its guaranteed that 2 out of 5 philosophers will eat at the SAME time
	for i := 0; i < 3; i++ {

		// CAS(Compare and Swap) operation. Need BOTH left and right forks available(means cur value is 0)
		if atomic.CompareAndSwapUint32(&p.left.pickFlag, 0, 1) && atomic.CompareAndSwapUint32(&p.right.pickFlag, 0, 1) {

			// we dont know which two philosophers will pick up the forks
			// totally depends on whoever picks up them fast ergo excruciatingly starving
			fmt.Printf("philosopher %d is eating\n", p.id)
			time.Sleep(2 * time.Second)

			// unset the flags
			atomic.StoreUint32(&p.left.pickFlag, 0)
			atomic.StoreUint32(&p.right.pickFlag, 0)

			fmt.Printf("philosopher %d is thinking\n", p.id)
		} else {
			// CAS failed. Put that goroutine to sleep using futex sys call
			// swap is equivalent to XChg
			//oldLeftValue := atomic.SwapUint32(&p.left.pickFlag, 2)
			//oldRightValue := atomic.SwapUint32(&p.right.pickFlag, 2)

			// futex sys call

		}

	}
}

func wineAndDineFUTEX() {
	// 5 philosophers and 5 forks

	forks := make([]*ForkFUTEX, 5)

	for i := 0; i < 5; i++ {
		forks[i] = &ForkFUTEX{
			pickFlag: 0, //suppose every fork is not initially available for adjacent philosopher to pick. They need to pick first.
		}
	}

	philosophers := make([]*PhilosopherFUTEX, 5)

	for i := 0; i < 5; i++ {
		philosophers[i] = &PhilosopherFUTEX{id: i, left: forks[i], right: forks[(i+1)%5]}
	}

	var wg sync.WaitGroup

	// start dining
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go philosophers[i].wantsToEatFUTEX(&wg)
	}

	wg.Wait()

}
