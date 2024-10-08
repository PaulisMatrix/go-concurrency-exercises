package main

import (
	"fmt"
	"sync"
	"time"
)

// https://leetcode.com/problems/the-dining-philosophers/description/
// this is a variant of the above confusing problem statement

type Fork struct {
	pick sync.Mutex
}

type Philosopher struct {
	id    int
	left  *Fork
	right *Fork
}

func (p Philosopher) wantsToEat(wg *sync.WaitGroup) {
	defer wg.Done()

	// its guaranteed that 2 out of 5 philosophers will eat at the SAME time
	for i := 0; i < 3; i++ {

		// try to pick up the forks
		// in C, this would be sem_wait(&p)
		// sem_wait(&p) is a decrementing operation
		/*
			func sem_wait(p *semStruct){
				p.semaState--
				if p.semaState < 0{
					// wait for someone else to release the lock
					// keep spinning
					// time.sleep(?)
				}
			}
		*/
		p.left.pick.Lock()
		p.right.pick.Lock()

		// we dont know which two philosophers will pick up the forks
		// totally depends on whoever picks up them fast ergo excruciatingly starving
		fmt.Printf("philosopher %d is eating\n", p.id)
		time.Sleep(2 * time.Second)

		// put down the forks
		// in C, this would be sem_post(&p)
		// sem_post(&p) is an incrementing operation.
		/*
			func sem_post(p *semStruct){
				p.semaState++ //thats it, its only job is to atomically increment the state and exit
			}

		*/
		p.left.pick.Unlock()
		p.right.pick.Unlock()

		fmt.Printf("philosopher %d is thinking\n", p.id)

	}
}

func wineAndDine() {
	// 5 philosophers and 5 forks

	forks := make([]*Fork, 5)

	for i := 0; i < 5; i++ {
		forks[i] = new(Fork)
	}

	philosophers := make([]*Philosopher, 5)

	for i := 0; i < 5; i++ {
		philosophers[i] = &Philosopher{id: i, left: forks[i], right: forks[(i+1)%5]}
	}

	var wg sync.WaitGroup

	// start dining
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go philosophers[i].wantsToEat(&wg)
	}

	wg.Wait()

}
