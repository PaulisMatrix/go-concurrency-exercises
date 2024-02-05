package main

import (
	"fmt"
	"sync"
)

func main() {
	// printInOrder()
	// alternateFooBar()
	// startFizzBuzzing()
	// startZeroEvenOdding()
	// wineAndDine()
	// wineAndDineCAS()
	buildH2O()
	// gotchas()

}

func gotchas() {
	// semaphores can be used for both, for locking/guarding a critical section (mutext.Lock() and mutext.Unlock())
	// as well as for ordering of process(sync.Waitgroup or chan). Both these pkgs use semaphores as synchornization primitive.

	// sync.Mutex docs says that: A locked mutex is not "associated" with particular goroutine
	// and it can be unlocked by an another goroutine as well.
	// this is correct since remember, its using semaphore in the back and semaphores are just signally mechanism.
	// there is nothing something inherently stopping other goroutines to access the critical section.
	// like goroutines being held back in shackles by something.
	// but semaphores just tell the current state(set, unset) of the critical section
	// more explanation here: https://www.reddit.com/r/golang/comments/1797dtu/comment/k54ckx2/

	// but do remember that unlocking a mutex which doesn't even have a lock in the first place will obv not work
	// and also recursive locking will not work since its basically hold and wait(already holding a lock and waiting to acquire a new one)
	// and that's one of the contenders for deadlocks.

	var wg sync.WaitGroup
	mutexChan := make(chan *sync.Mutex)

	wg.Add(1)
	go func() {
		// create a lock and pass the lock in the chan
		var mu sync.Mutex
		fmt.Println("goroutine1 taking the lock")

		mu.Lock()
		mutexChan <- &mu

		fmt.Println("goroutine2 will release the lock not held by it technically")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		mu := <-mutexChan
		mu.Unlock()

		fmt.Println("releasing the lock originally held by goroutine1")
		wg.Done()
	}()

	wg.Wait()

}
