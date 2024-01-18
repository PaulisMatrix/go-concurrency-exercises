package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// https://blog.stackademic.com/mutex-internals-in-golang-1624749f35a6

// problem with this CAS operation is, the philosopher goroutines which wont be able to do a successful
// CAS operation will just wait burning CPU cycles?

type ForkCAS struct {
	// set, unset this flag for CAS operationn
	// set = 1 = Availble for pickup
	// unset = 0 = Not Available for pickup
	pickFlag uint32
}

type PhilosopherCAS struct {
	id    int
	left  *ForkCAS
	right *ForkCAS
}

func (p PhilosopherCAS) wantsToEatCAS(wg *sync.WaitGroup) {
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
			// this philosopher routine's CAS fails since some other one has acquired it right now.
			// by the time the other philosopher unset's the flag, this routine has already finished
			// so it never gets to eat
			// try adding a sleep here?

			fmt.Printf("philosopher %d failed CAS!?\n", p.id)
			time.Sleep(2 * time.Second)

			// even if sleep is added. this is not deterministic since we don't know whether
			// this philosopher routine is even awake to set the flag when the other one has unset it already.
		}

	}
}

func wineAndDineCAS() {
	// 5 philosophers and 5 forks

	forks := make([]*ForkCAS, 5)

	for i := 0; i < 5; i++ {
		forks[i] = &ForkCAS{
			pickFlag: 0, //suppose every fork is not initially available for adjacent philosopher to pick. They need to pick first.
		}
	}

	philosophers := make([]*PhilosopherCAS, 5)

	for i := 0; i < 5; i++ {
		philosophers[i] = &PhilosopherCAS{id: i, left: forks[i], right: forks[(i+1)%5]}
	}

	var wg sync.WaitGroup

	// start dining
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go philosophers[i].wantsToEatCAS(&wg)
	}

	wg.Wait()

}
