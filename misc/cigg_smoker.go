package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// https://w3.cs.jmu.edu/kirkpams/OpenCSF/Books/csf/html/CigSmokers.html

type Smoke struct {
	paperTobaccoBuffer chan string
	paperMatchBuffer   chan string
	tobaccoMatchBuffer chan string

	done chan bool
}

func (s *Smoke) agent(wg *sync.WaitGroup) {
	// agent randomly produces two of the three items at a time
	defer wg.Done()
loop:
	for {
		select {
		case <-s.done:
			// close all the channels first so each smoker waiting on the channel buffer exits its own loop
			close(s.paperMatchBuffer)
			close(s.paperTobaccoBuffer)
			close(s.tobaccoMatchBuffer)
			break loop
		default:
			number := rand.Intn(10) % 3
			switch number {
			case 0:
				// produce paper and tobacoo
				// this is like sema_post(), incrementing/adding to the semaphore/buffer resp.
				s.paperTobaccoBuffer <- "paper"
				s.paperTobaccoBuffer <- "tobacco"
			case 1:
				// produce paper and match
				s.paperMatchBuffer <- "paper"
				s.paperMatchBuffer <- "match"
			case 2:
				// produce tobacco and match
				s.tobaccoMatchBuffer <- "tobacoo"
				s.tobaccoMatchBuffer <- "match"
			}
		}
	}
}

func (s *Smoke) tobacco_smoker(wg *sync.WaitGroup) {
	// tobacco smoker has infinite tobacco
	// it just need a paper and a match to smoke
	defer wg.Done()

	// retrieve from the buffer
	// for now assume we dont know which item is this
	// but its guaranteed to be either paper or match
	// so if its added to the buffer, pull from it
	// this will get blocked until both are available which is what we want

	// this is like sema_wait(), decrementing/consuming from the semaphore/buffer resp.
	i := 1
	for item := range s.paperMatchBuffer {
		if i%2 == 0 {
			fmt.Println("tobacoo smoker finished smoking.waiting for few secs")
			time.Sleep(1 * time.Second)
		}
		fmt.Print("got: ", item)
		i++
	}

}

func (s *Smoke) paper_smoker(wg *sync.WaitGroup) {
	// paper smoker has infinite paper
	// it just need a tobacco and a match to smoke
	defer wg.Done()

	i := 1
	for item := range s.tobaccoMatchBuffer {
		if i%2 == 0 {
			fmt.Println("paper smoker finished smoking.waiting for few secs")
			time.Sleep(1 * time.Second)
		}
		fmt.Print("got: ", item)
		i++
	}

}

func (s *Smoke) match_smoker(wg *sync.WaitGroup) {
	// match smoker has infinite match
	// it just need a paper and a tobacco to smoke
	defer wg.Done()

	i := 1
	for item := range s.paperTobaccoBuffer {
		if i%2 == 0 {
			fmt.Println("match smoker finished smoking.waiting for few secs")
			time.Sleep(1 * time.Second)
		}
		fmt.Print("got: ", item)
		i++
	}
}

func ciggSmoker() {
	var wg sync.WaitGroup
	s := Smoke{
		paperTobaccoBuffer: make(chan string, 2),
		paperMatchBuffer:   make(chan string, 2),
		tobaccoMatchBuffer: make(chan string, 2),

		done: make(chan bool),
	}

	timeout := time.After(2 * time.Second)

	wg.Add(4)
	go s.agent(&wg)
	go s.match_smoker(&wg)
	go s.paper_smoker(&wg)
	go s.tobacco_smoker(&wg)

	<-timeout
	fmt.Println("sending done. waiting for all smokers to exit")
	// send one
	s.done <- true

	wg.Wait()

}
